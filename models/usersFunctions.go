package models

import (
	"errors"

	"github.com/astaxie/beego/orm"
)

func (m *User) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *User) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *User) InsertStudent(p *Profile) error {
	m.Profile = p
	o := orm.NewOrm()
	err := o.Begin()
	if _, err = o.Insert(p); err != nil {
		o.Rollback()
		return err
	}
	if _, err = o.Insert(m); err != nil {
		o.Rollback()
		return err
	}
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

func (m *User) InsertTeacher(p *Profile, t *Teacher, tg *TeacherTags) error {
	m.Profile = p
	p.Teacher = t
	t.TeacherTags = tg
	o := orm.NewOrm()
	err := o.Begin()
	if _, err = o.Insert(tg); err != nil {
		o.Rollback()
		return err
	}
	if _, err = o.Insert(t); err != nil {
		o.Rollback()
		return err
	}
	if _, err = o.Insert(p); err != nil {
		o.Rollback()
		return err
	}
	if _, err = o.Insert(m); err != nil {
		o.Rollback()
		return err
	}

	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

func (m *User) LoadProfile() error {
	if _, err := orm.NewOrm().LoadRelated(m, "Profile"); err != nil {
		return err
	}
	return nil
}

func (p *Profile) LoadTeacher() error {
	if _, err := orm.NewOrm().LoadRelated(p, "Teacher"); err != nil {
		return err
	}
	return nil
}

func (d *UserData) Insert() error {
	if _, err := orm.NewOrm().Insert(d); err != nil {
		return err
	}
	return nil
}

func VerifyEmail(emailData string) error {
	var m UserData
	o := orm.NewOrm()
	err := o.QueryTable(m).Filter("type", "emailVerify").Filter("data", emailData).One(&m)
	if err == orm.ErrMultiRows {
		// Have multiple records
		return errors.New("Returned Multi Rows Not One")
	}
	if err == orm.ErrNoRows {
		// No result
		return errors.New("Not row found")
	}
	user := m.User
	user.IsActive = true
	err = user.Update("IsActive")
	if err != nil {
		return err
	}
	return nil
}

func (m *UserData) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Profile) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (t *Teacher) Read(fields ...string) error {
	if err := orm.NewOrm().Read(t, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Teacher) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (t *Teacher) LoadProfile() error {
	if _, err := orm.NewOrm().LoadRelated(t, "Profile"); err != nil {
		return err
	}
	return nil
}

func Users() orm.QuerySeter {
	var table User
	return orm.NewOrm().QueryTable(table).OrderBy("-Id")
}

func Teachers() orm.QuerySeter {
	var table Teacher
	return orm.NewOrm().QueryTable(table).OrderBy("IsActive", "-Id")
}

func GetTeachers() (t []*Teacher, num int64, err error) {
	teachers := Teachers()
	num, err = teachers.All(&t)
	return
}

func VerifyResume(id int)  error {
	teacher := Teacher{Id: id}
	if err := teacher.Read("Id"); err != nil {
		return err
	}
	teacher.IsActive = true
	err := teacher.Update("IsActive")
	if err != nil {
		return err
	}
	return nil
}