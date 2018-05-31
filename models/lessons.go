package models

import (
	"github.com/astaxie/beego/orm"
)

// User 使用者帳號密碼
type LessonSchedule struct {
	Id      int
	Week    int      `orm:"default(1)"`
	H0      int      `orm:"default(-1)"`
	H1      int      `orm:"default(-1)"`
	H2      int      `orm:"default(-1)"`
	H3      int      `orm:"default(-1)"`
	H4      int      `orm:"default(-1)"`
	H5      int      `orm:"default(-1)"`
	H6      int      `orm:"default(-1)"`
	H7      int      `orm:"default(-1)"`
	H8      int      `orm:"default(-1)"`
	H9      int      `orm:"default(-1)"`
	H10     int      `orm:"default(-1)"`
	H11     int      `orm:"default(-1)"`
	H12     int      `orm:"default(-1)"`
	H13     int      `orm:"default(-1)"`
	H14     int      `orm:"default(-1)"`
	H15     int      `orm:"default(-1)"`
	H16     int      `orm:"default(-1)"`
	H17     int      `orm:"default(-1)"`
	H18     int      `orm:"default(-1)"`
	H19     int      `orm:"default(-1)"`
	H20     int      `orm:"default(-1)"`
	H21     int      `orm:"default(-1)"`
	H22     int      `orm:"default(-1)"`
	H23     int      `orm:"default(-1)"`
	Profile *Profile `orm:"rel(fk)"` // RelForeignKey relation
}

// multiple fields unique key
func (u *LessonSchedule) TableUnique() [][]string {
	return [][]string{
		[]string{"Week", "Profile"},
	}
}

func init() {
	// Need to register model in init
	orm.RegisterModel(new(LessonSchedule))

}
