package controllers

import (
	"WebPartice/models"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"html/template"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

func Authenticate(email string, password string) (user *models.User, err error) {
	msg := "invalid email or password."
	user = &models.User{Email: email}

	if err := user.Read("Email"); err != nil {
		if err.Error() == "<QuerySeter> no row found" {
			err = errors.New(msg)
		}
		return user, err
	} else if user.Id < 1 {
		// No user
		return user, errors.New(msg)
	} else if user.Password != password {
		// No matched password
		return user, errors.New(msg)
	} else {
		user.Lastlogintime = time.Now()
		user.Update("Lastlogintime")
		return user, nil
	}
}

func SignupStudent(u *models.User, p *models.Profile) (int, error) {
	var (
		err error
		msg string
	)

	if models.Users().Filter("email", u.Email).Exist() {
		msg = "was already regsitered input email address."
		return 0, errors.New(msg)
	}

	h := sha256.New()
	h.Write([]byte(u.Email))
	h.Write([]byte(u.Password))
	u.Password = string(base64.URLEncoding.EncodeToString(h.Sum(nil)))

	err = u.Insert(p)
	if err != nil {
		return 0, err
	}

	return u.Id, err
}

func IsValid(model interface{}) (err error) {
	valid := validation.Validation{}
	b, err := valid.Valid(model)
	if !b {
		for _, err := range valid.Errors {
			beego.Warning(err.Key, ":", err.Message)
			return errors.New(err.Message)
			// return errors.New(fmt.Sprintf("%s: %s", err.Key, err.Message))
		}
	}
	return nil
}

// AuthController 使用者控制器，含是否登陸資訊
type AuthController struct {
	beego.Controller
	Userinfo *models.User
	IsLogin  bool
}

func (c *AuthController) Get() {
	c.Prepare()
	if c.IsLogin {
		c.Ctx.Redirect(302, c.URLFor("UsersController.Index"))
		return
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "login.html"
}

func (c *AuthController) SetLogin(user *models.User) {
	c.SetSession("userinfo", user.Id)
}

func (c *AuthController) DelLogin() {
	c.DestroySession()
}

func (c *AuthController) Prepare() {
	c.IsLogin = c.GetSession("userinfo") != nil
	if c.IsLogin {
		c.Userinfo = c.GetLogin()
	}
}

func (c *AuthController) GetLogin() *models.User {
	u := &models.User{Id: c.GetSession("userinfo").(int)}
	return u
}

func (c *AuthController) Login() {
	c.Prepare()
	if c.IsLogin {
		c.Ctx.Redirect(302, c.URLFor("UsersController.Index"))
		return
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "login.html"
	flash := beego.NewFlash()
	email := c.GetString("Email")
	password := c.GetString("Password")

	h := sha256.New()
	h.Write([]byte(email))
	h.Write([]byte(password))

	user, err := Authenticate(email, string(base64.URLEncoding.EncodeToString(h.Sum(nil))))
	if err != nil || user.Id < 1 {
		flash.Warning(err.Error())
		flash.Store(&c.Controller)
		return
	}

	flash.Success("Success logged in")
	flash.Store(&c.Controller)

	c.SetLogin(user)

	c.Redirect(c.URLFor("IndexController.Get"), 303)
}

func (c *AuthController) Logout() {
	c.DelLogin()
	flash := beego.NewFlash()
	flash.Success("Success logged out")
	flash.Store(&c.Controller)

	c.Ctx.Redirect(302, c.URLFor("LoginController.Login"))
}

func (c *AuthController) Signup() {
	c.TplName = "signup.html"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())

	if !c.Ctx.Input.IsPost() {
		return
	}

	var err error
	flash := beego.NewFlash()

	u := &models.User{}
	if err = c.ParseForm(u); err != nil {
		flash.Error("Signup invalid!")
		flash.Store(&c.Controller)
		return
	}
	p := &models.Profile{}
	if err = c.ParseForm(p); err != nil {
		flash.Error("Signup invalid!")
		flash.Store(&c.Controller)
		return
	}
	if err = IsValid(u); err != nil {
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		return
	}

	identity := c.Input()["Identity"]
	if len(identity) == 0 {
		flash.Warning("Please choose your identity.")
		flash.Store(&c.Controller)
		return
	}
	p.Identity = identity[0]
	if identity[0] == "student" {
		id, err := SignupStudent(u, p)
		if err != nil || id < 1 {
			flash.Warning(err.Error())
			flash.Store(&c.Controller)
			return
		}
	} else if identity[0] == "teacher" {
		// TODO register teacher
		flash.Warning("Not yet complete")
		flash.Store(&c.Controller)
		return
	} else {
		flash.Warning("Signup invalid!")
		flash.Store(&c.Controller)
		return
	}

	flash.Success("Register user")
	flash.Store(&c.Controller)

	c.SetLogin(u)

	c.Redirect(c.URLFor("UsersController.Index"), 303)
}
