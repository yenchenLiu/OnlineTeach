package controllers

import (
	"WebPartice/models"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

func (c *BaseController) Load() {
	IsLogin := c.GetSession("userinfo") != nil
	c.Data["IsLogin"] = IsLogin
	if IsLogin {
		user := &models.User{Id: c.GetSession("userinfo").(int)}
		user.Read()
		if user.IsActive == false {
			flash := beego.NewFlash()
			flash.Warning("The user hasn't verified yet. Please check the e-mail to verify your account.")
			flash.Store(&c.Controller)
			c.Abort("403")
			return
		}

		user.LoadProfile()
		if c.GetSession("isTeacher") == nil {
			IsTeacher := user.Profile.Identity == "teacher"
			c.SetSession("isTeacher", IsTeacher)
		}
		c.Data["Name"] = user.Profile.Name
		c.Data["Identity"] = user.Profile.Identity
	}
}
