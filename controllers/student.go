package controllers

import (
	"OnlineTeach/lib"
	"OnlineTeach/models"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type StudentAuditingController struct {
	BaseController
}

func (c *StudentAuditingController) Prepare() {
	if c.GetSession("IsStudent") != true {
		c.Abort("401")
	}
	c.LoadSession()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
}

func (c *StudentAuditingController) Get() {
	student := models.Student{Id: c.GetSession("student").(int)}
	student.Read("Id")
	student.LoadAuditing()

	var auditing []map[string]string
	for _, item := range student.StudentAuditings {
		if item.Teacher != nil {
			item.LoadTeacher()
			item.Teacher.LoadProfile()
			auditing = append(auditing, map[string]string{
				"Id": strconv.Itoa(item.Id), "Day": getWeekChinese(int64(item.Day)),
				"Hour": strconv.Itoa(item.Hour), "Status": item.Status,
				"TeacherId": strconv.Itoa(item.Teacher.Id), "TeacherName": item.Teacher.Profile.Name,
				"Skype":      item.Teacher.Profile.Skype,
				"LessonDate": item.ArrangeDate.Format("2006-01-02")})
		} else {
			auditing = append(auditing, map[string]string{
				"Id": strconv.Itoa(item.Id), "Day": getWeekChinese(int64(item.Day)),
				"Hour": strconv.Itoa(item.Hour), "Status": item.Status,
				"TeacherId": "", "TeacherName": "", "Skype": "", "LessonDate": item.ArrangeDate.Format("2006-01-02")})
		}
	}
	c.Data["auditing"] = auditing

	c.TplName = "student/auditing.html"
}

func (c *StudentAuditingController) Post() {
	student := models.Student{Id: c.GetSession("student").(int)}
	student.Read("Id")

	o := orm.NewOrm()
	qs := o.QueryTable("StudentAuditing")
	cnt, err := qs.Filter("student__id", student.Id).Filter("status", "安排中").Count()
	if err != nil {
		fmt.Println(err)
		c.Get()
		return
	}
	if cnt > 0 {
		flash := beego.NewFlash()
		flash.Warning("目前已有安排的試聽，請結束後再重新安排。")
		flash.Store(&c.Controller)
		c.Get()
		return
	}

	weekday, _ := strconv.ParseInt(c.Input()["day"][0], 10, 64)
	hour, _ := strconv.ParseInt(c.Input()["hour"][0], 10, 64)
	auditing := new(models.StudentAuditing)
	auditing.Student = &student
	auditing.Day = int(weekday)
	auditing.Hour = int(hour)
	auditing.Status = "安排中"
	d, _ := time.ParseDuration("24h")
	t := time.Now().Add(2 * d)
	fmt.Println(t)
	// (today weekday - expected weekday)
	days := 7 - int((t.Weekday() - time.Weekday(weekday)))
	if days >= 7 {
		days -= 7
	}
	fmt.Println(int(t.Weekday()))
	fmt.Println(int(time.Weekday(weekday)))
	fmt.Println(int(days))

	t = t.Add(time.Duration(days) * d)
	fmt.Println(t)
	auditing.ArrangeDate = t

	if err := auditing.Insert(); err != nil {
		fmt.Println(err)
	}

	c.Get()
}

func getWeekChinese(week int64) string {
	switch week {
	case 0:
		return "星期日"
	case 1:
		return "星期一"
	case 2:
		return "星期二"
	case 3:
		return "星期三"
	case 4:
		return "星期四"
	case 5:
		return "星期五"
	case 6:
		return "星期六"
	}
	return ""
}

type StudentDepositController struct {
	BaseController
}

func (c *StudentDepositController) Prepare() {
	if c.GetSession("IsStudent") != true {
		c.Abort("401")
	}
	c.LoadSession()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
}

func (c *StudentDepositController) Get() {
	c.TplName = "student/deposit.html"
}

func (c *StudentDepositController) Post() {
	money := 0
	value := c.Input()["deposit"][0]
	switch value {
	case "10":
		money = 1000
	case "20":
		money = 2000
	case "40":
		money = 4000
	default:
		c.Get()
		return
	}
	data := strings.Split(lib.PayMoney(c.GetSession("student").(int), money), "&")
	post := make(map[string]string)

	for _, item := range data {
		s := strings.Split(item, "=")
		post[s[0]] = s[1]
	}
	fmt.Println(post)
	c.Data["Post"] = post

	c.TplName = "student/redirectToEZpay.tpl"
}
