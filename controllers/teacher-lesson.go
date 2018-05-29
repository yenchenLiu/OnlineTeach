package controllers

type LessonController struct {
	BaseController
}

func (c *LessonController) Get() {
	c.Load()
	if c.GetSession("isTeacher") != true {
		c.Abort("401")
	}
	var lesson [24][7]bool
	for index := 7; index <= 18; index++ {
		lesson[index][0] = true
		lesson[index][1] = true
		lesson[index][2] = true
		lesson[index][3] = true
	}
	c.Data["lesson"] = lesson
	c.TplName = "teacher/lesson_schedule.html"
}
