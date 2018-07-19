package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type RatingRecords struct {
	Id                 int
	TeachingQuality    int8
	ExplainClearly     int8
	TeachingAttraction int8
	Description        string    `orm:"type(text)"`
	Created            time.Time `orm:"auto_now_add;type(datetime)"`
	Updated            time.Time `orm:"auto_now;type(datetime)"`
	Student            *Student  `orm:"rel(fk)"` // RelForeignKey relation
	Teacher            *Teacher  `orm:"rel(fk)"` // RelForeignKey relation
}

// multiple fields unique key
func (r *RatingRecords) TableUnique() [][]string {
	return [][]string{
		[]string{"Student", "Teacher"},
	}
}

func init() {
	// Need to register model in init
	orm.RegisterModel(new(RatingRecords))
}
