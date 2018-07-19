package controllers

import (
	"time"
	"OnlineTeach/models"
	"html/template"
	"strconv"

	"github.com/astaxie/beego/orm"
)

type AdminReviewResumeController struct {
	BaseController
}

func (c *AdminReviewResumeController) Prepare() {
	if c.GetSession("IsAdmin") != true {
		c.Abort("401")
	}
	c.LoadSession()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
}

func (c *AdminReviewResumeController) Get() {
	var teachersData []map[string]string

	teachers, _, _ := models.GetTeachers()
	for _, v := range teachers {
		v.LoadProfile()
		teachersData = append(teachersData, map[string]string{
			"Id":          strconv.Itoa(v.Id),
			"Name":        v.Profile.Name,
			"Teaching":    strconv.Itoa(v.TeachingYears),
			"Proficiency": v.Proficiency, "Resume": v.Resume,
			"Youtube":  v.Youtube,
			"IsActive": strconv.FormatBool(v.IsActive)})
	}

	c.Data["teachers"] = teachersData
	c.TplName = "admin/reviewResume.html"
}

func (c *AdminReviewResumeController) Post() {
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

type AdminProcessWithdrawMoney struct {
	BaseController
}

func (w *AdminProcessWithdrawMoney) Prepare() {
	if w.GetSession("IsAdmin") != true {
		w.Abort("401")
	}
	w.LoadSession()
	w.Data["xsrfdata"] = template.HTML(w.XSRFFormHTML())
}

func (w *AdminProcessWithdrawMoney) Get() {

	var withdrawRecord []models.WithdrawRecord
	o := orm.NewOrm()
	// 取得老師提領紀錄
	qs := o.QueryTable("WithdrawRecord")
	if _, err := qs.All(&withdrawRecord); err != nil {
		return
	}
	location, _ := time.LoadLocation("Asia/Taipei")
	var withdrawRecords []map[string]string
	for _, item := range withdrawRecord {
		item.LoadProfile()
		local := item.Created
		local = local.In(location)
		withdrawRecords = append(withdrawRecords, map[string]string{
			"profileId": strconv.Itoa(item.Profile.Id),
			"Name":      item.Profile.Name,
			"Id":        strconv.Itoa(item.Id),
			"Points":    strconv.FormatFloat(item.Points, 'f', 2, 64),
			"PayPal":    item.PayPal,
			"Process":   item.Process,
			"Description": item.Description,
			"Created":   local.Format("2006-01-02 15:04:05")})
	}
	w.Data["withdraw"] = withdrawRecords

	w.TplName = "admin/teacherWithdraw.html"
}
