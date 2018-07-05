package routers

import (
	"OnlineTeach/controllers"
	"OnlineTeach/scheduling"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.IndexController{})
	beego.Router("/index", &controllers.IndexController{})
	beego.Router("/login", &controllers.AuthController{}, "get:Get;post:Login")
	beego.Router("/logout", &controllers.AuthController{}, "get:Logout")
	beego.Router("/signup", &controllers.AuthController{}, "*:Signup")
	beego.Router("/verify/:verify", &controllers.AuthController{}, "get:VerifyEmail")
	beego.Router("/signup/teacher", &controllers.AuthController{}, "*:SignupTeacher")
	beego.Router("/student/auditing", &controllers.StudentAuditingController{})
	beego.Router("/student/deposit", &controllers.StudentDepositController{})
	beego.Router("/student/new/lesson", &controllers.NewLessonController{})
	beego.Router("/student/teacherInformation/:Id", &controllers.TeacherInformation{})
	beego.Router("/student/courseList", &controllers.CourseListForStudent{})
	beego.Router("/teacher/lesson", &controllers.LessonController{})
	beego.Router("/teacher/auditing", &controllers.TeacherAuditing{})
	beego.Router("/teacher/courseList", &controllers.CourseListForTeacher{})
	beego.Router("/admin/reviewresume", &controllers.AdminReviewResumeController{})
	beego.Router("/admin/reviewresume/:name/:file", &controllers.AdminReviewResumeController{}, "get:Download")
	
	// 電子商務
	beego.Router("/ecpay/receive", &controllers.ECPayController{})

	// 刷新資料
	beego.Router("/scheduling/course", &scheduling.CourseController{})
}
