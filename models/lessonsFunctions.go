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