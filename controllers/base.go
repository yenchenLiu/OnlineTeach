package controllers

import (
	"WebPartice/models"
	"fmt"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

func (c *BaseController) Load() {
	IsLogin := c.GetSession("userinfo") != nil
	fmt.Print(IsLogin, "\n")
	c.Data["IsLogin"] = IsLogin
	if IsLogin {
		fmt.Print(c.GetSession("userinfo").(int), "\n")
		user := &models.User{Id: c.GetSession("userinfo").(int)}
		if err := user.Read("Id"); err != nil {
			fmt.Printf("ERR: %v\n", err)
		}
		user.LoadProfile()
		if c.GetSession("isTeacher") == nil {
			IsTeacher := user.Profile.Identity == "teacher"
			c.SetSession("isTeacher", IsTeacher)
		}
		c.Data["Name"] = user.Profile.FirstName
		c.Data["Identity"] = user.Profile.Identity
	}
}
