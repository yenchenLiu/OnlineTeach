package routers

import (
	"WebPartice/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.IndexController{})
	beego.Router("/index", &controllers.IndexController{})
	beego.Router("/login", &controllers.AuthController{}, "get:Get;post:Login")
	beego.Router("/logout", &controllers.AuthController{}, "get:Logout")
	beego.Router("/signup", &controllers.AuthController{}, "*:Signup")
	beego.Router("/signup/teacher", &controllers.AuthController{}, "*:SignupTeacher")
	beego.Router("/teacher/lesson", &controllers.LessonController{},)
}
