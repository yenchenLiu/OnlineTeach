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

// StudentDepositController 學生儲值
type StudentDepositController struct {
	BaseController
}

func (s *StudentDepositController) Prepare() {
	if s.GetSession("IsStudent") != true {
		s.Abort("401")
	}
	s.LoadSession()
	s.Data["xsrfdata"] = template.HTML(s.XSRFFormHTML())
}

func (s *StudentDepositController) Get() {
	profile := models.Profile{Id: s.GetSession("ProfileId").(int)}

	s.Data["Points"] = strconv.FormatFloat(profile.Points, 'f', 1, 64)
	s.TplName = "student/deposit.html"
}

func (s *StudentDepositController) Post() {
	money := 0
	value := s.Input()["deposit"][0]
	switch value {
	case "10":
		money = 1000
	case "20":
		money = 2000
	case "40":
		money = 4000
	default:
		s.Get()
		return
	}
	data := strings.Split(lib.PayMoney(s.GetSession("student").(int), money), "&")
	post := make(map[string]string)

	for _, item := range data {
		s := strings.Split(item, "=")
		post[s[0]] = s[1]
	}
	fmt.Println(post)
	s.Data["Post"] = post
	s.Data["ecpayurl"] = beego.AppConfig.String("ECPAYUrl")

	s.TplName = "student/redirectToEZpay.tpl"
}

type NewLessonController struct {
	BaseController
}

func (n *NewLessonController) Prepare() {
	if n.GetSession("IsStudent") != true {
		n.Abort("401")
	}
	n.LoadSession()
	n.Data["xsrfdata"] = template.HTML(n.XSRFFormHTML())
}

func (n *NewLessonController) Get() {
	var teachersData []map[string]string

	teachers, _, _ := models.GetTeachers()
	for _, v := range teachers {
		v.LoadProfile()
		if v.IsActive == true {
			teachersData = append(teachersData, map[string]string{
				"Id":             strconv.Itoa(v.Id),
				"Name":           v.Profile.Name,
				"Rating":         strconv.FormatFloat(v.AverageRating, 'f', 1, 64),
				"TotalClassHour": strconv.FormatFloat(v.TotalClassHour, 'f', 1, 64)})
		}
	}

	n.Data["teachers"] = teachersData
	n.TplName = "student/newLesson.html"
}

type TeacherInformation struct {
	BaseController
}

func (t *TeacherInformation) Prepare() {
	if t.GetSession("IsStudent") != true {
		t.Abort("401")
	}
	t.LoadSession()
	t.Data["xsrfdata"] = template.HTML(t.XSRFFormHTML())
}

func (t *TeacherInformation) Get() {
	Id, _ := strconv.ParseInt(t.GetString(":Id"), 10, 64)
	teacher := models.Teacher{Id: int(Id)}
	if err := teacher.Read("Id"); err != nil {
		t.Abort("404")
	}
	teacher.LoadProfile()

	schedules := models.LoadSchedule(teacher.Profile.Id)
	var lessons [18]map[int]int
	for index := 0; index < 18; index++ {
		lessons[index] = make(map[int]int)
	}
	// 6-23
	for index, schedule := range schedules {
		lessons[0][index] = schedule.H6
		lessons[1][index] = schedule.H7
		lessons[2][index] = schedule.H8
		lessons[3][index] = schedule.H9
		lessons[4][index] = schedule.H10
		lessons[5][index] = schedule.H11
		lessons[6][index] = schedule.H12
		lessons[7][index] = schedule.H13
		lessons[8][index] = schedule.H14
		lessons[9][index] = schedule.H15
		lessons[10][index] = schedule.H16
		lessons[11][index] = schedule.H17
		lessons[12][index] = schedule.H18
		lessons[13][index] = schedule.H19
		lessons[14][index] = schedule.H20
		lessons[15][index] = schedule.H21
		lessons[16][index] = schedule.H22
		lessons[17][index] = schedule.H23
	}
	t.Data["lessons"] = lessons

	t.TplName = "student/teacherInformation.html"
}
