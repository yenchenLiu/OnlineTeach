package controllers

type IndexController struct {
	BaseController
}

func (c *IndexController) Get() {
	c.TplName = "index.html"
	c.LoadSession()
	if c.GetSession("IsTeacher") == true {
		c.TplName = "teacher/dashboard.html"
	}

	if c.GetSession("IsStudent") == true {
		c.TplName = "student/dashboard.html"
	}

	if c.GetSession("IsAdmin") == true {
		c.TplName = "admin/dashboard.html"
	}
}
