package controllers

type IndexController struct {
	BaseController
}

func (c *IndexController) Get() {
	c.Load()
	c.TplName = "index.html"
	
	if c.GetSession("isTeacher") == true {
		c.TplName = "teacher/dashboard.html"
	}
}
