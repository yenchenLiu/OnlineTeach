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
	Profile       *Profile  `orm:"rel(one)";` // OneToOne relation
}

// Profile 使用者個人資料
type Profile struct {
	Id        int
	Name      string
	Identity  string
	Points    float64           `orm:"digits(12);decimals(2);default(0.00)"`
	User      *User             `orm:"reverse(one);on_delete(set_null)"` // Reverse relationship (optional)
	Teacher   *Teacher          `orm:"null;rel(one)"`                    // Reverse relationship (optional)
	Schedules []*LessonSchedule `orm:"reverse(many)"`                    // reverse relationship of fk
}

type Teacher struct {
	Id            int
	TeachingYears int
	Proficiency   string `orm:"type(text)"`
	Resume        string
	IsActive      bool         `orm:"default(false)"`
	Profile       *Profile     `orm:"reverse(one)"` // Reverse relationship (optional)
	TeacherTags   *TeacherTags `orm:"rel(one)"`     // OneToOne relation
}

type TeacherTags struct {
	Id       int
	Child    bool     `orm:"default(false)"`
	Beginner bool     `orm:"default(false)"`
	Advanced bool     `orm:"default(false)"`
	TOEIC    bool     `orm:"default(false)"`
	TOFEL    bool     `orm:"default(false)"`
	Teacher  *Teacher `orm:"reverse(one)"` // Reverse relationship (optional)
}

func init() {
	// Need to register model in init
	orm.RegisterModel(new(User), new(Profile), new(Teacher), new(TeacherTags))

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

func (m *User) InsertStudent(p *Profile) error {
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

func (m *User) InsertTeacher(p *Profile, t *Teacher, tg *TeacherTags) error {
	m.Profile = p
	p.Teacher = t
	t.TeacherTags = tg
	o := orm.NewOrm()
	err := o.Begin()
	o.Insert(tg)
	o.Insert(t)
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
