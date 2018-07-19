package scheduling

import (
	"OnlineTeach/models"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type jsonResponse struct {
	Data string `json:"data"`
}
type CourseController struct {
	beego.Controller
}

func (c *CourseController) Prepare() {
	c.EnableXSRF = false
}

func (c *CourseController) Post() {
	if len(c.Input()["secretKey"]) != 1 {
		c.Abort("404")
	}
	if c.Input()["secretKey"][0] != beego.AppConfig.String("schedulingKey") {
		c.Abort("404")
	}
	var responseData string
	responseData = "success"
	RefreshALLCourse()

	response := jsonResponse{Data: responseData}
	c.Data["json"] = response
	c.ServeJSON()
}

func RefreshALLCourse() {

	courses, _, _ := models.GetAllActiveRegistration()
	for _, course := range courses {
		o := orm.NewOrm()
		if cnt, err := o.QueryTable("courseRecord").Filter("Status", "即將上課").Filter("CourseRegistration", course).Count(); err == nil {

			fmt.Println(cnt)
			if cnt != 0 {
				continue
			}

			d, _ := time.ParseDuration("24h")
			t := time.Now().Add(d)
			// (today weekday - expected weekday)
			weekday := course.ClassWeek
			fmt.Println(weekday)
			days := 7 - int((t.Weekday() - time.Weekday(weekday)))
			if days >= 7 {
				days -= 7
			}
			t = t.Add(time.Duration(days) * d)
			fmt.Println(t)
			// 沒有課程單，新增課程單
			courseRecord := new(models.CourseRecord)
			courseRecord.Status = "即將上課"
			courseRecord.ClassTimeDate = t
			courseRecord.ClassTimeHour = course.ClassHour
			courseRecord.CourseRegistration = course
			courseRecord.Insert()

		} else {
			fmt.Println(err)
		}

	}

}
