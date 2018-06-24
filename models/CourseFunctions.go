package models

import (
	"github.com/astaxie/beego/orm"
)

func (l *CourseSchedule) Read(fields ...string) error {
	if err := orm.NewOrm().Read(l, fields...); err != nil {
		return err
	}
	return nil
}

func (l *CourseSchedule) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(l, fields...); err != nil {
		return err
	}
	return nil
}

func (l *CourseSchedule) Insert() error {
	if _, err := orm.NewOrm().Insert(l); err != nil {
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

func GetCourseRegistrationFromStudent(student *Student) (c []*CourseRegistration, num int64, err error) {
	var table CourseRegistration
	courses := orm.NewOrm().QueryTable(table).Filter("Student", student).Filter("IsActive", true).OrderBy("-Id")
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
