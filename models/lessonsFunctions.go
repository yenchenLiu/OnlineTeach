package models

import (
	"github.com/astaxie/beego/orm"
)

func (l *LessonSchedule) Read(fields ...string) error {
	if err := orm.NewOrm().Read(l, fields...); err != nil {
		return err
	}
	return nil
}

func (l *LessonSchedule) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(l, fields...); err != nil {
		return err
	}
	return nil
}

func (l *LessonSchedule) Insert() error {
	if _, err := orm.NewOrm().Insert(l); err != nil {
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

func LoadSchedule(Id int) []LessonSchedule {
	var schedules []LessonSchedule
	o := orm.NewOrm()
	qs := o.QueryTable("LessonSchedule")
	qs.Filter("Profile", Id).All(&schedules)
	return schedules
}

func UpdateSchedule(TeacherId int, week int, hour int, value int) error {
	return nil
}
