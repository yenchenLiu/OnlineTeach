package controllers

import (
	"WebPartice/models"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

func (c *BaseController) LoadSession() {
	c.Data["IsLogin"] = c.GetSession("IsLogin")
	c.Data["IsStudent"] = c.GetSession("IsStudent")
	c.Data["IsTeacher"] = c.GetSession("IsTeacher")
	c.Data["IsAdmin"] = c.GetSession("IsAdmin")
	c.Data["Name"] = c.GetSession("Name")
	c.Data["Identity"] = c.GetSession("Identity")
	c.Data["Uri"] = c.Ctx.Input.URI()

}

func (c *BaseController) Load() {
	IsLogin := c.GetSession("userinfo") != nil
	c.SetSession("IsLogin", IsLogin)
	if IsLogin {
		user := &models.User{Id: c.GetSession("userinfo").(int)}
		user.Read()
		if user.IsActive == false {
			flash := beego.NewFlash()
			flash.Warning("The user hasn't verified yet. Please check the e-mail to verify your account.")
			flash.Store(&c.Controller)
			c.DestroySession()
			c.Abort("403")
			return
		}

		user.LoadProfile()
		c.SetSession("ProfileId", user.Profile.Id)
		if c.GetSession("IsTeacher") == nil {
			IsTeacher := user.Profile.Identity == "teacher"
			c.SetSession("IsTeacher", IsTeacher)
			if IsTeacher == true {
				user.Profile.LoadTeacher()
				if IsTeacherActive := user.Profile.Teacher.IsActive; IsTeacherActive == false {
					flash := beego.NewFlash()
					flash.Warning("The resume is under review and will be notified by email when the audit is approved.")
					flash.Store(&c.Controller)
					c.DestroySession()
					c.Abort("403")
					return
				}
				c.SetSession("teacher", user.Profile.Teacher.Id)
			}
		}
		if c.GetSession("IsStudent") == nil {
			IsStudent := user.Profile.Identity == "student"
			c.SetSession("IsStudent", IsStudent)
			user.Profile.LoadStudent()
			if IsStudent == true {
				c.SetSession("student", user.Profile.Student.Id)
			}
		}

		if c.GetSession("IsAdmin") == nil {
			IsTeacher := user.Profile.Identity == "admin"
			c.SetSession("IsAdmin", IsTeacher)
		}
		c.SetSession("Name", user.Profile.Name)
		c.SetSession("Identity", user.Profile.Identity)
	}
}
