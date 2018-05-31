package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// User 使用者帳號密碼
type User struct {
	Id            int
	Email         string
	Password      string
	Registertime  time.Time `orm:"auto_now_add;type(datetime)"`
	Lastlogintime time.Time `orm:"type(datetime);null" form:"-"`
	IsActive      bool      `orm:"default(false)"`
	Profile       *Profile  `orm:"rel(one)";on_delete(set_null)` // OneToOne relation
}

// Profile 使用者個人資料
type Profile struct {
	Id        int
	Name      string
	Identity  string
	User      *User             `orm:"reverse(one)"`  // Reverse relationship (optional)
	Teacher   *Teacher          `orm:"null;rel(one)"` // Reverse relationship (optional)
	Schedules []*LessonSchedule `orm:"reverse(many)"` // reverse relationship of fk
}

type Teacher struct {
	Id            int
	TeachingYears int
	Proficiency   string `orm:"type(text)"`

	IsActive bool     `orm:"default(false)"`
	Profile  *Profile `orm:"reverse(one)"` // Reverse relationship (optional)
}

func init() {
	// Need to register model in init
	orm.RegisterModel(new(User), new(Profile), new(Teacher))

}

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

func (m *User) Insert(p *Profile) error {
	m.Profile = p
	o := orm.NewOrm()
	err := o.Begin()
	o.Insert(p)
	o.Insert(m)
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

func (m *Profile) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func Users() orm.QuerySeter {
	var table User
	return orm.NewOrm().QueryTable(table).OrderBy("-Id")
}
