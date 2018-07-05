package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CourseSchedule struct {
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
func (u *CourseSchedule) TableUnique() [][]string {
	return [][]string{
		[]string{"Week", "Profile"},
	}
}

type CourseRegistration struct {
	Id        int
	ClassWeek int8
	ClassHour int8
	IsActive  bool      `orm:"default(true)"`
	Points    float64   `orm:"digits(12);decimals(2);default(0.00)"`
	Teacher   *Teacher  `orm:"rel(fk)"` // RelForeignKey relation
	Student   *Student  `orm:"rel(fk)"` // RelForeignKey relation
	Created   time.Time `orm:"auto_now_add;type(datetime)"`
	Updated   time.Time `orm:"auto_now;type(datetime)"`
}

type CourseRecord struct {
	Id                 int
	Status             string    // 即將上課、上課中、已結束 。 上課中到結束時須轉移點數
	ClassTimeDate      time.Time `orm:"type(date)"`
	ClassTimeHour      int8
	TeachingMaterial   string              `orm:"default()"`
	TeachingDetail     string              `orm:"type(text);default()"`
	CourseRegistration *CourseRegistration `orm:"rel(fk)"` // RelForeignKey relation
	Created            time.Time           `orm:"auto_now_add;type(datetime)"`
	Updated            time.Time           `orm:"auto_now;type(datetime)"`
}

type StudentAuditing struct {
	Id          int
	Day         int
	Hour        int
	Status      string
	ArrangeDate time.Time `orm:"type(date)"`
	Created     time.Time `orm:"auto_now_add;type(datetime)"`
	Updated     time.Time `orm:"auto_now;type(datetime)"`
	Student     *Student  `orm:"rel(fk)"`      // RelForeignKey relation
	Teacher     *Teacher  `orm:"rel(fk);null"` // RelForeignKey relation
}

func init() {
	// Need to register model in init
	orm.RegisterModel(new(CourseSchedule), new(CourseRegistration), new(CourseRecord), new(StudentAuditing))

}
