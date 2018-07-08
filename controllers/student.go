package controllers

import (
	"OnlineTeach/lib"
	"OnlineTeach/models"
	"fmt"
	"html/template"
	"net/url"
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
	// (today weekday - expected weekday)
	days := 7 - int((t.Weekday() - time.Weekday(weekday)))
	if days >= 7 {
		days -= 7
	}

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
	teacherData := make(map[string]string)
	youtubeUrl, err := url.Parse(teacher.Youtube)
	if err != nil {
		teacherData["Youtube"] = ""
	} else {
		if len(youtubeUrl.Query()["v"]) == 1 {
			teacherData["Youtube"] = youtubeUrl.Query()["v"][0]
		}
	}
	teacherData["Name"] = teacher.Profile.Name
	teacherData["AverageRating"] = strconv.FormatFloat(teacher.AverageRating, 'f', 1, 64)
	teacherData["TotalClassHour"] = strconv.FormatFloat(teacher.TotalClassHour, 'f', 1, 64)
	teacherData["Proficiency"] = teacher.Proficiency
	t.Data["lessons"] = lessons
	t.Data["teacherData"] = teacherData
	t.TplName = "student/teacherInformation.html"
}

func (t *TeacherInformation) Post() {
	Id, _ := strconv.ParseInt(t.GetString(":Id"), 10, 64)
	teacher := models.Teacher{Id: int(Id)}
	if err := teacher.Read("Id"); err != nil {
		t.Abort("404")
	}
	teacher.LoadProfile()
	profile := models.Profile{Id: t.GetSession("ProfileId").(int)}
	profile.LoadStudent()
	student := profile.Student
	value := strings.Split(t.Input()["addLesson"][0], "_")
	var schedule models.CourseSchedule
	week, _ := strconv.ParseInt(value[0], 10, 64)
	hour, _ := strconv.ParseInt(value[1], 10, 64)
	flash := beego.NewFlash()

	o := orm.NewOrm()

	// 檢查點數是否足夠選課
	if profile.Points < 2 {
		flash.Warning("點數不足，請先至課程儲值購買點數")
		flash.Store(&t.Controller)
		t.Get()
		return
	}

	// 檢查與更改學生課表
	if err := o.QueryTable("CourseSchedule").Filter("Profile", t.GetSession("ProfileId").(int)).Filter("Week", week).One(&schedule); err != nil {
		t.Redirect(t.URLFor("TeacherInformation.Get"), 302)
	}

	switch getField(&schedule, "H"+value[1]) {
	case -1: // 成功選取
		courseRegistration := new(models.CourseRegistration)
		courseRegistration.ClassWeek = int8(week)
		courseRegistration.ClassHour = int8(hour)
		courseRegistration.Points = teacher.ClassValue
		courseRegistration.Student = student
		courseRegistration.IsActive = true
		courseRegistration.Teacher = &teacher
		courseRegistration.Insert()

		setField(&schedule, "H"+value[1], courseRegistration.Id)
		schedule.Update("H" + value[1])

		var scheduleTeacher models.CourseSchedule
		// 更改老師課表
		if err := o.QueryTable("CourseSchedule").Filter("Profile", teacher.Profile.Id).Filter("Week", week).One(&scheduleTeacher); err != nil {
			t.Redirect(t.URLFor("TeacherInformation.Get"), 302)
		}
		switch getField(&scheduleTeacher, "H"+value[1]) {
		case 0:
			setField(&scheduleTeacher, "H"+value[1], courseRegistration.Id)
			scheduleTeacher.Update("H" + value[1])
		default: // 非法操作
			flash.Warning("此時段已安排其他課程")
			flash.Store(&t.Controller)
			t.Get()
			return
		}

	default: // 已有課程 return error
		flash.Warning("此時段已安排其他課程")
		flash.Store(&t.Controller)
		t.Get()
		return
	}

	// TODO 產生課程單

	flash.Success("選課成功，以下為選課老師資訊")
	flash.Store(&t.Controller)
	t.Redirect(t.URLFor("CourseListForStudent.Get"), 302)
}

type CourseListForStudent struct {
	BaseController
}

func (c *CourseListForStudent) Prepare() {
	if c.GetSession("IsStudent") != true {
		c.Abort("401")
	}
	c.LoadSession()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
}

func (c *CourseListForStudent) Get() {
	student := models.Student{Id: c.GetSession("student").(int)}
	student.Read("Id")
	courses, num, _ := models.GetCourseRegistrationFromStudent(&student)
	if num == 0 {
		// TODO 沒有選課時，顯示教學畫面
	} else {
		var courseDatas []map[string]string
		for _, course := range courses {
			course.LoadTeacher()
			course.Teacher.LoadProfile()
			courseRecord := models.CourseRecord{Status: "即將上課", CourseRegistration: course}
			t := map[string]string{
				"ID":           strconv.Itoa(course.Id),
				"Point":        strconv.FormatFloat(course.Points, 'f', 2, 64),
				"teacherID":    strconv.Itoa(course.Teacher.Id),
				"teacherName":  course.Teacher.Profile.Name,
				"teacherSkype": course.Teacher.Profile.Skype}

			if err := courseRecord.Read("Status", "CourseRegistration"); err == nil {
				t["classTime"] = courseRecord.ClassTimeDate.Format("2006-01-02") + " " + strconv.Itoa(int(courseRecord.ClassTimeHour)) + ":00"
			} else {
				fmt.Println(err)
			}
			courseDatas = append(courseDatas, t)
		}
		c.Data["courseDatas"] = courseDatas

	}
	c.TplName = "student/courseList.html"
}
func (c *CourseListForStudent) Post() {
	if len(c.Input()["cancelLesson"]) != 1 {
		c.Abort("400")
	}

	lessonId, _ := strconv.ParseInt(c.Input()["cancelLesson"][0], 10, 64)
	courseRegistration := models.CourseRegistration{Id: int(lessonId)}
	courseRegistration.Read("Id")
	// 判斷是否為本人取消課程
	courseRegistration.LoadStudent()
	if courseRegistration.Student.Id != c.GetSession("student").(int) {
		c.Abort("401")
	}
	courseRecord := models.CourseRecord{Status: "即將上課", CourseRegistration: &courseRegistration}
	if err := courseRecord.Read("Status", "CourseRegistration"); err == nil {
		courseRecord.Status = "已取消"
		courseRecord.Update("Status")
	} else {
		fmt.Println(err)
	}
	courseRegistration.IsActive = false
	courseRegistration.Update("IsActive")
	// 取消老師與學生課程
	courseRegistration.LoadTeacher()
	courseRegistration.LoadStudent()
	courseRegistration.Student.LoadProfile()
	courseRegistration.Teacher.LoadProfile()

	var student_schedule models.CourseSchedule
	o := orm.NewOrm()
	week := courseRegistration.ClassWeek
	hour := courseRegistration.ClassHour
	// 檢查與更改學生課表
	if err := o.QueryTable("CourseSchedule").Filter("Profile", courseRegistration.Student.Profile).Filter("Week", week).One(&student_schedule); err != nil {
		fmt.Println(err)
	}

	setField(&student_schedule, "H"+strconv.Itoa(int(hour)), -1)
	student_schedule.Update("H" + strconv.Itoa(int(hour)))

	var teacher_schedule models.CourseSchedule
	if err := o.QueryTable("CourseSchedule").Filter("Profile", courseRegistration.Teacher.Profile).Filter("Week", week).One(&teacher_schedule); err != nil {
		fmt.Println(err)
	}
	setField(&teacher_schedule, "H"+strconv.Itoa(int(hour)), 0)
	teacher_schedule.Update("H" + strconv.Itoa(int(hour)))

	c.Get()
}
