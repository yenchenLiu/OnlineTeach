package controllers

import (
	"html/template"
)

type StudentAuditingController struct {
	BaseController
}

func (c *StudentAuditingController) Get() {
	if c.GetSession("IsStudent") != true {
		c.Abort("401")
	}
	c.LoadSession()

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())

	c.TplName = "student/auditing.html"
}
