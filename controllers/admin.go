package controllers

import (
	"WebPartice/models"
	"fmt"
	"html/template"
	"strconv"
)

type AdminReviewResumeController struct {
	BaseController
}

func (c *AdminReviewResumeController) Get() {
	if c.GetSession("IsAdmin") != true {
		c.Abort("401")
	}
	c.LoadSession()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	var teachersData []map[string]string

	teachers, num, err := models.GetTeachers()
	fmt.Print(teachers, "\n")
	fmt.Print(num, "\n")
	fmt.Print(err, "\n")
	for _, v := range teachers {
		v.LoadProfile()
		teachersData = append(teachersData, map[string]string{
			"Id":          strconv.Itoa(v.Id),
			"Name":        v.Profile.Name,
			"Teaching":    strconv.Itoa(v.TeachingYears),
			"Proficiency": v.Proficiency, "Resume": v.Resume,
			"IsActive": strconv.FormatBool(v.IsActive)})
	}

	c.Data["teachers"] = teachersData
	c.TplName = "admin/reviewResume.html"
}

func (c *AdminReviewResumeController) Post() {
	if c.GetSession("IsAdmin") != true {
		c.Abort("401")
	}

	for _, v := range c.Input()["Teacher[]"] {
		id, _ := strconv.ParseInt(v, 10, 64)
		models.VerifyResume(int(id))
	}
	c.Redirect(c.URLFor("AdminReviewResumeController.Get"), 302)
}

func (c *AdminReviewResumeController) Download() {
	if c.GetSession("IsAdmin") != true {
		c.Abort("401")
	}
	c.Ctx.Output.Download("./resumes/"+c.GetString(":file"), c.GetString(":name")+".pdf")
}