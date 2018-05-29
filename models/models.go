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
	Lastlogintime time.Time `orm:"type(datetime);null" form:"-"`
	Profile       *Profile  `orm:"rel(one)"` // OneToOne relation
}

// Profile 使用者個人資料
type Profile struct {
	Id   int
	Name string
	Age  int16
	User *User `orm:"reverse(one)"` // Reverse relationship (optional)
}

func init() {
	// Need to register model in init
	orm.RegisterModel(new(User), new(Profile))

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
