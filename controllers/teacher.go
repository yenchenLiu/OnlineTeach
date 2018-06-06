package controllers

import (
	"OnlineTeach/models"
	"fmt"
	"html/template"
	"reflect"
	"strconv"
	"strings"

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
	schedules := models.LoadSchedule(c.GetSession("userinfo").(int))
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

	var schedule models.LessonSchedule
	var responseData string

	week, _ := strconv.ParseInt(value[0], 10, 64)
	// hour, _ := strconv.ParseInt(value[1], 10, 64)
	o := orm.NewOrm()
	if err := o.QueryTable("LessonSchedule").Filter("Profile", c.GetSession("userinfo").(int)).Filter("Week", week).One(&schedule); err != nil {
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

func getField(v *models.LessonSchedule, field string) int {
	r := reflect.ValueOf(v)
	fmt.Println(r)
	f := reflect.Indirect(r).FieldByName(field)
	fmt.Println(f)
	return int(f.Int())
}

func setField(v *models.LessonSchedule, field string, data int) {
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
	qs := o.QueryTable("StudentAuditing")

	if _, err := qs.Filter("teacher__id", teacher.Id).All(&auditings); err != nil {
		fmt.Println(err)
	}
	var teacherAuditing []map[string]string
	for _, item := range auditings {
		item.LoadStudent()
		item.Student.LoadProfile()
		teacherAuditing = append(teacherAuditing, map[string]string{
			"StudentName": item.Student.Profile.Name, "Skype": item.Student.Profile.Skype,
			"LessonDate": item.ArrangeDate.Format("2006-01-02"),
			"Hour":       strconv.Itoa(item.Hour)})
		students = append(students, item.Student.Id)
	}
	this.Data["teacherAuditing"] = teacherAuditing
	var auditings_2 []models.StudentAuditing
	qs = o.QueryTable("StudentAuditing")
	if _, err := qs.Filter("Status", "安排中").Exclude("student__id__in", students).All(&auditings_2); err != nil {
		fmt.Println(err)
	}
	fmt.Println(auditings_2)

	var auditing []map[string]string
	for _, item := range auditings_2 {
		var schedule models.LessonSchedule
		o.QueryTable("LessonSchedule").Filter("profile__id", this.GetSession("ProfileId").(int)).Filter("Week", item.Day).One(&schedule)
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
