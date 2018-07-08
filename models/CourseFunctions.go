package models

import (
	"github.com/astaxie/beego/orm"
)

func (c *CourseSchedule) Read(fields ...string) error {
	if err := orm.NewOrm().Read(c, fields...); err != nil {
		return err
	}
	return nil
}

func (c *CourseSchedule) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(c, fields...); err != nil {
		return err
	}
	return nil
}

func (c *CourseSchedule) Insert() error {
	if _, err := orm.NewOrm().Insert(c); err != nil {
		return err
	}
	return nil
}


func (c *CourseRegistration) Read(fields ...string) error {
	if err := orm.NewOrm().Read(c, fields...); err != nil {
		return err
	}
	return nil
}

func (c *CourseRegistration) Insert() error {
	if _, err := orm.NewOrm().Insert(c); err != nil {
		return err
	}
	return nil
}


func (c *CourseRegistration) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(c, fields...); err != nil {
		return err
	}
	return nil
}


func (c *CourseRecord) Read(fields ...string) error {
	if err := orm.NewOrm().Read(c, fields...); err != nil {
		return err
	}
	return nil
}

func (c *CourseRecord) Insert() error {
	if _, err := orm.NewOrm().Insert(c); err != nil {
		return err
	}
	return nil
}

func (c *CourseRecord) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(c, fields...); err != nil {
		return err
	}
	return nil
}


func (c *CourseRecord) LoadRegistration() error {
	if _, err := orm.NewOrm().LoadRelated(c, "CourseRegistration"); err != nil {
		return err
	}
	return nil
}

func GetCourseRegistrationFromStudent(student *Student) (c []*CourseRegistration, num int64, err error) {
	var table CourseRegistration
	courses := orm.NewOrm().QueryTable(table).Filter("Student", student).Filter("IsActive", true).OrderBy("-Id")
	num, err = courses.All(&c)
	return
}

func GetCourseRegistrationFromTeacher(teacher *Teacher) (c []*CourseRegistration, num int64, err error) {
	var table CourseRegistration
	courses := orm.NewOrm().QueryTable(table).Filter("Teacher", teacher).Filter("IsActive", true).OrderBy("-Id")
	num, err = courses.All(&c)
	return
}

func GetAllActiveRegistration() (c []*CourseRegistration, num int64, err error) {
	var table CourseRegistration
	courses := orm.NewOrm().QueryTable(table).Filter("IsActive", true).OrderBy("-Id")
	num, err = courses.All(&c)
	return
}

func (c *CourseRegistration) LoadTeacher() error {
	if _, err := orm.NewOrm().LoadRelated(c, "Teacher"); err != nil {
		return err
	}
	return nil
}

func (c *CourseRegistration) LoadStudent() error {
	if _, err := orm.NewOrm().LoadRelated(c, "Student"); err != nil {
		return err
	}
	return nil
}

func (s *StudentAuditing) Read(fields ...string) error {
	if err := orm.NewOrm().Read(s, fields...); err != nil {
		return err
	}
	return nil
}

func (s *StudentAuditing) LoadStudent() error {
	if _, err := orm.NewOrm().LoadRelated(s, "Student"); err != nil {
		return err
	}
	return nil
}

func (s *StudentAuditing) LoadTeacher() error {
	if _, err := orm.NewOrm().LoadRelated(s, "Teacher"); err != nil {
		return err
	}
	return nil
}

func (s *StudentAuditing) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(s, fields...); err != nil {
		return err
	}
	return nil
}

func (a *StudentAuditing) Insert() error {
	if _, err := orm.NewOrm().Insert(a); err != nil {
		return err
	}
	return nil
}

func LoadSchedule(Id int) []CourseSchedule {
	var schedules []CourseSchedule
	o := orm.NewOrm()
	qs := o.QueryTable("CourseSchedule")
	qs.Filter("Profile", Id).All(&schedules)
	return schedules
}

func UpdateSchedule(TeacherId int, week int, hour int, value int) error {
	return nil
}
