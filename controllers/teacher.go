package controllers

import (
	"time"
	"OnlineTeach/models"
	"fmt"
	"html/template"
	"reflect"
	"strconv"
	"strings"

	"github.com/astaxie/beego"

	"github.com/astaxie/beego/orm"
)

type LessonController struct {
	BaseController
}

func (c *LessonController) Prepare() {
	if c.GetSession("IsTeacher") != true {
		c.Abort("401")
	}
	c.LoadSession()
	c.EnableXSRF = false

}

func (c *LessonController) Get() {
	schedules := models.LoadSchedule(c.GetSession("ProfileId").(int))
	fmt.Println(schedules)
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
	c.Data["lessons"] = lessons
	c.TplName = "teacher/lesson_schedule.html"
}

type jsonResponse struct {
	Data string `json:"data"`
}

func (c *LessonController) Post() {
	value := strings.Split(c.Input()["lesson_checkbox"][0], "_")

	var schedule models.CourseSchedule
	var responseData string

	week, _ := strconv.ParseInt(value[0], 10, 64)
	// hour, _ := strconv.ParseInt(value[1], 10, 64)
	o := orm.NewOrm()
	if err := o.QueryTable("CourseSchedule").Filter("Profile", c.GetSession("ProfileId").(int)).Filter("Week", week).One(&schedule); err != nil {
		c.Redirect(c.URLFor("LessonController.Get"), 302)
	}

	switch getField(&schedule, "H"+value[1]) {
	case -1:
		responseData = "Open " + getWeek(week) + " " + value[1] + " o'clock"
		setField(&schedule, "H"+value[1], 0)
	case 0:
		responseData = "Close " + getWeek(week) + " " + value[1] + " o'clock"
		setField(&schedule, "H"+value[1], -1)
	}
	schedule.Update("H" + value[1])
	fmt.Println(getField(&schedule, "H"+value[1]))
	response := jsonResponse{Data: responseData}
	c.Data["json"] = response
	c.ServeJSON()
}

func getField(v *models.CourseSchedule, field string) int {
	r := reflect.ValueOf(v)
	fmt.Println(r)
	f := reflect.Indirect(r).FieldByName(field)
	fmt.Println(f)
	return int(f.Int())
}

func setField(v *models.CourseSchedule, field string, data int) {
	r := reflect.ValueOf(v)
	fmt.Println(r)
	reflect.Indirect(r).FieldByName(field).SetInt(int64(data))
}

func getWeek(week int64) string {
	switch week {
	case 0:
		return "Sunday"
	case 1:
		return "Monday"
	case 2:
		return "Tuesday"
	case 3:
		return "Wensday"
	case 4:
		return "Thursday"
	case 5:
		return "Friday"
	case 6:
		return "Saturday"
	}
	return ""
}

type TeacherAuditing struct {
	BaseController
}

func (this *TeacherAuditing) Prepare() {
	if this.GetSession("IsTeacher") != true {
		this.Abort("401")
	}
	this.LoadSession()
	this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
}

func (this *TeacherAuditing) Get() {
	teacher := models.Teacher{Id: this.GetSession("teacher").(int)}
	teacher.Read("Id")
	students := make([]int, 0)
	var auditings []models.StudentAuditing
	o := orm.NewOrm()

	// 找出老師所有教的學生
	qs := o.QueryTable("StudentAuditing")
	if _, err := qs.Filter("teacher__id", teacher.Id).All(&auditings); err != nil {
		fmt.Println(err)
	}
	var teacherAuditing []map[string]string
	for _, item := range auditings {
		item.LoadStudent()
		item.Student.LoadProfile()
		fmt.Println(item.Student.Id)
		teacherAuditing = append(teacherAuditing, map[string]string{
			"StudentName": item.Student.Profile.Name, "Skype": item.Student.Profile.Skype,
			"LessonDate": item.ArrangeDate.Format("2006-01-02"),
			"Hour":       strconv.Itoa(item.Hour)})
		students = append(students, item.Student.Id)
	}
	this.Data["teacherAuditing"] = teacherAuditing
	var auditings2 []models.StudentAuditing
	qs = o.QueryTable("StudentAuditing")
	// 如果老師教的學生數不等於0
	if len(students) != 0 {
		if _, err := qs.Filter("Status", "安排中").Exclude("student__id__in", students).All(&auditings2); err != nil {
			fmt.Println(err)
		}
	} else {
		if _, err := qs.Filter("Status", "安排中").All(&auditings2); err != nil {
			fmt.Println(err)
		}
	}

	var auditing []map[string]string
	for _, item := range auditings2 {
		var schedule models.CourseSchedule
		o.QueryTable("CourseSchedule").Filter("profile__id", this.GetSession("ProfileId").(int)).Filter("Week", item.Day).One(&schedule)
		if getField(&schedule, "H"+strconv.Itoa(item.Hour)) == 0 {
			item.LoadStudent()
			item.Student.LoadProfile()
			// TODO 課程日期產生程式
			auditing = append(auditing, map[string]string{
				"Id":          strconv.Itoa(item.Id),
				"StudentName": item.Student.Profile.Name, "Skype": item.Student.Profile.Skype,
				"LessonDate": item.ArrangeDate.Format("2006-01-02"),
				"Hour":       strconv.Itoa(item.Hour)})
		}
	}
	this.Data["auditing"] = auditing

	this.TplName = "teacher/auditing.html"
}

func (this *TeacherAuditing) Post() {
	id, _ := strconv.ParseInt(this.Input()["AuditingId"][0], 10, 64)
	auditing := models.StudentAuditing{Id: int(id)}
	if auditing.Teacher == nil {
		teacher := models.Teacher{Id: this.GetSession("teacher").(int)}
		teacher.Read("Id")
		auditing.Teacher = &teacher
		auditing.Status = "預約完成"
		auditing.Update("Teacher", "Status")
	}

	this.Get()
}

type CourseListForTeacher struct {
	BaseController
}

func (c *CourseListForTeacher) Prepare() {
	if c.GetSession("IsTeacher") != true {
		c.Abort("401")
	}
	c.LoadSession()
}

func (c *CourseListForTeacher) Get() {
	teacher := models.Teacher{Id: c.GetSession("teacher").(int)}
	teacher.Read("Id")
	courses, num, _ := models.GetCourseRegistrationFromTeacher(&teacher)
	if num == 0 {
		// TODO 沒有選課時，顯示教學畫面
	} else {
		var courseDatas []map[string]string
		for _, course := range courses {
			course.LoadStudent()
			course.Student.LoadProfile()
			courseRecord := models.CourseRecord{Status: "即將上課", CourseRegistration: course}
			t := map[string]string{
				"ID":           strconv.Itoa(course.Id),
				"Point":        strconv.FormatFloat(course.Points, 'f', 2, 64),
				"studentID":    strconv.Itoa(course.Student.Id),
				"studentName":  course.Student.Profile.Name,
				"studentSkype": course.Student.Profile.Skype}
			if err := courseRecord.Read("Status", "CourseRegistration"); err == nil {
				t["classTime"] = courseRecord.ClassTimeDate.Format("2006-01-02") + " " + strconv.Itoa(int(courseRecord.ClassTimeHour)) + ":00"
			}
			courseDatas = append(courseDatas, t)
		}
		c.Data["courseDatas"] = courseDatas

	}
	c.TplName = "teacher/courseList.html"
}

type WithdrawMoney struct {
	BaseController
}

func (w *WithdrawMoney) Prepare() {
	if w.GetSession("IsTeacher") != true {
		w.Abort("401")
	}
	w.LoadSession()
	w.Data["xsrfdata"] = template.HTML(w.XSRFFormHTML())
}

func (w *WithdrawMoney) Get() {
	profile := models.Profile{Id: w.GetSession("ProfileId").(int)}
	profile.Read("Id")
	profile.LoadTeacher()
	w.Data["Points"] = strconv.FormatFloat(profile.Points, 'f', 1, 64)
	w.Data["PointsFloat"] = float64(profile.Points)
	if len(profile.Teacher.PayPal) > 22 {
		w.Data["PayPal"] = profile.Teacher.PayPal
	} else {
		w.Data["PayPal"] = "https://www.paypal.me/"
	}

	var withdrawRecord []models.WithdrawRecord
	o := orm.NewOrm()
	// 取得取款紀錄
	qs := o.QueryTable("WithdrawRecord")
	if _, err := qs.Filter("profile__id", profile.Id).All(&withdrawRecord); err != nil {
		fmt.Println(err)
		return
	}
	location, _ := time.LoadLocation("Asia/Taipei")
	var withdrawRecords []map[string]string
	for _, item := range withdrawRecord {
		local := item.Created
		local = local.In(location)
		withdrawRecords = append(withdrawRecords, map[string]string{
			"Points":  strconv.FormatFloat(item.Points, 'f', 2, 64),
			"PayPal":  item.PayPal,
			"Process": item.Process,
			"Created": local.Format("2006-01-02 15:04:05")})
	}
	w.Data["withdraw"] = withdrawRecords

	w.TplName = "teacher/withdraw.html"
}

func (w *WithdrawMoney) Post() {

	point, _ := strconv.ParseInt(w.Input()["withdraw"][0], 10, 64)
	paypal := w.Input()["paypal"][0]
	if len(paypal) < 24 {
		flash := beego.NewFlash()
		flash.Warning("PayPal URL error, please fill in the URL of paypal.me. Please see the note for details.")
		flash.Store(&w.Controller)
		w.Get()
		return
	}

	profile := models.Profile{Id: w.GetSession("ProfileId").(int)}
	profile.Read("Id")

	// 點數不足，跳出錯誤頁面
	if float64(point) > profile.Points {
		w.Abort("400")
	}
	profile.LoadTeacher()
	profile.Teacher.PayPal = paypal
	profile.Teacher.Update("PayPal")

	withdrawUser := models.User{Email: "withdraw@daychen.tw"}
	withdrawUser.Read("Email")
	withdrawUser.LoadProfile()

	money := strconv.Itoa(int(point * 85))

	o := orm.NewOrm()
	o.Using("default")

	err := o.Begin()
	withdrawRecord := new(models.WithdrawRecord)
	withdrawRecord.Points = float64(point)
	withdrawRecord.PayPal = paypal
	withdrawRecord.Profile = &profile
	withdrawRecord.Process = "suspending"
	withdrawRecord.Description = profile.Name + "提領" + w.Input()["withdraw"][0] + "點，要匯款" + money + "元"

	pointsTrade := new(models.PointsTrade)
	pointsTrade.Points = float64(point)
	pointsTrade.Description = profile.Name + "提領" + w.Input()["withdraw"][0] + "點，要匯款" + money + "元"
	pointsTrade.ProfileGiver = &profile
	pointsTrade.ProfileReceiver = withdrawUser.Profile

	profile.Points -= float64(point)
	withdrawUser.Profile.Points += float64(point)

	if err == nil {
		_, err = o.Insert(withdrawRecord)
	}
	if err == nil {
		_, err = o.Insert(pointsTrade)
	}
	if err == nil {
		_, err = o.Update(&profile, "Points")
	}
	if err == nil {
		_, err = o.Update(withdrawUser.Profile, "Points")
	}
	if err != nil {
		o.Rollback()
	} else {
		o.Commit()
	}

	w.Get()
}
