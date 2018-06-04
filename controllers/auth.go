package controllers

import (
	"WebPartice/lib"
	"WebPartice/models"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"html/template"
	"net/url"
	"path"
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

func SignupVerify(u *models.User, p *models.Profile) error {
	var (
		err error
		msg string
	)

	if models.Users().Filter("email", u.Email).Exist() {
		msg = "was already regsitered input email address."
		return errors.New(msg)
	}

	return err
}

func SignupTeacher(u *models.User, p *models.Profile, t *models.Teacher, tg *models.TeacherTags) (int, error) {
	var (
		err error
	)

	h := sha256.New()
	h.Write([]byte(u.Email))
	h.Write([]byte(u.Password))
	u.Password = string(base64.URLEncoding.EncodeToString(h.Sum(nil)))

	err = u.InsertTeacher(p, t, tg)
	if err != nil {
		return 0, err
	}

	userData := new(models.UserData)
	emailVerify := string(base64.URLEncoding.EncodeToString(h.Sum(nil)))
	userData.User = u
	userData.Type = "emailVerify"
	userData.Data = emailVerify
	userData.Insert()

	lib.SendVerifyMail("s412172010@gmail.com", emailVerify)
	return u.Id, err
}

func SignupStudent(u *models.User, p *models.Profile) (int, error) {
	var (
		err error
	)

	h := sha256.New()
	h.Write([]byte(u.Email))
	h.Write([]byte(u.Password))
	u.Password = string(base64.URLEncoding.EncodeToString(h.Sum(nil)))

	h = sha256.New()
	h.Write([]byte(u.Email))
	h.Write([]byte(time.Now().String()))
	err = u.InsertStudent(p)
	if err != nil {
		return 0, err
	}

	userData := new(models.UserData)
	emailVerify := string(base64.URLEncoding.EncodeToString(h.Sum(nil)))
	userData.User = u
	userData.Type = "emailVerify"
	userData.Data = emailVerify
	userData.Insert()
	lib.SendVerifyMail("s412172010@gmail.com", emailVerify)
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
	BaseController
	Userinfo *models.User
	IsLogin  bool
}

func (c *AuthController) Get() {
	c.Prepare()
	if c.IsLogin {
		c.Ctx.Redirect(302, c.URLFor("IndexController.Get"))
		return
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "login.html"
}

func (c *AuthController) VerifyEmail() {
	flash := beego.NewFlash()

	if err := models.VerifyEmail(c.GetString(":verify")); err != nil {
		flash.Warning(err.Error())
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, c.URLFor("AuthController.Login"))
		return
	}

	flash.Notice("Success verify your account, please login again.")
	flash.Store(&c.Controller)
	c.DelSession("userinfo")
	c.Ctx.Redirect(302, c.URLFor("AuthController.Login"))
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
		c.Ctx.Redirect(302, c.URLFor("IndexController.Get"))
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
	c.Load()
	c.Redirect(c.URLFor("IndexController.Get"), 303)
}

func (c *AuthController) Logout() {
	c.DelLogin()
	flash := beego.NewFlash()
	flash.Success("Success logged out")
	flash.Store(&c.Controller)

	c.Ctx.Redirect(302, c.URLFor("AuthController.Login"))
}

func (c *AuthController) Signup() {
	c.TplName = "signup.html"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())

	if !c.Ctx.Input.IsPost() {
		if c.GetSession("RegisterTeacher") != nil {
			c.Redirect(c.URLFor("AuthController.SignupTeacher"), 303)
		}
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
		err := SignupVerify(u, p)
		if err != nil {
			flash.Warning(err.Error())
			flash.Store(&c.Controller)
			return
		}
		id, err := SignupStudent(u, p)
		if err != nil || id < 1 {
			flash.Warning(err.Error())
			flash.Store(&c.Controller)
			return
		}
	} else if identity[0] == "teacher" {

		// TODO 注意是否是傳指標
		c.SetSession("RegisterTeacher", c.Input())
		err := SignupVerify(u, p)
		if err != nil {
			c.DestroySession()
			flash.Warning(err.Error())
			flash.Store(&c.Controller)
			return
		}
		c.Redirect(c.URLFor("AuthController.SignupTeacher"), 303)
		return

	} else {
		flash.Warning("Signup invalid!")
		flash.Store(&c.Controller)
		return
	}

	flash.Success("Register user")
	flash.Store(&c.Controller)

	c.SetLogin(u)

	c.Redirect(c.URLFor("IndexController.Get"), 303)
}

func (c *AuthController) SignupTeacher() {
	if c.GetSession("RegisterTeacher") == nil {
		c.Redirect(c.URLFor("IndexController.Get"), 303)
	}
	c.TplName = "teacher/signup.html"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	if !c.Ctx.Input.IsPost() {
		return
	}
	flash := beego.NewFlash()
	_, header, err := c.GetFile("ResumeFile")
	if err != nil {
		flash.Warning(err.Error())
		flash.Store(&c.Controller)
		return
	}
	if path.Ext(header.Filename) != ".pdf" {
		flash.Warning("Please upload pdf file.")
		flash.Store(&c.Controller)
		return
	}

	register := c.GetSession("RegisterTeacher").(url.Values)

	user := new(models.User)
	profile := new(models.Profile)
	tg := new(models.TeacherTags)
	t := &models.Teacher{}
	if err = c.ParseForm(t); err != nil {
		flash.Error("Signup invalid!")
		flash.Store(&c.Controller)
		return
	}
	for _, v := range c.Input()["TeacherTags[]"] {
		switch v {
		case "Child":
			tg.Child = true
			fallthrough
		case "Beginner":
			tg.Beginner = true
		case "Advanced":
			tg.Advanced = true
		case "TOEIC":
			tg.TOEIC = true
		case "TOFEL":
			tg.TOFEL = true
		}
	}

	profile.Identity = register["Identity"][0]
	profile.Name = register["Name"][0]
	profile.Skype = register["Skype"][0]
	user.Email = register["Email"][0]
	user.Password = register["Password"][0]
	b := []byte(register["Email"][0])
	has := md5.Sum(b)
	filename := hex.EncodeToString(has[:])
	t.Resume = filename + ".pdf"
	if err := c.SaveToFile("ResumeFile", "./resumes/"+filename+".pdf"); err != nil {
		flash.Warning(err.Error())
		flash.Store(&c.Controller)
		return
	}

	id, err := SignupTeacher(user, profile, t, tg)
	if err != nil || id < 1 {
		flash.Warning(err.Error())
		flash.Store(&c.Controller)
		return
	}

	flash.Success("Register user")
	flash.Store(&c.Controller)
	c.DestroySession()
	c.Redirect(c.URLFor("IndexController.Get"), 303)
}
