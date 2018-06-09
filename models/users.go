package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// User 使用者帳號密碼
type User struct {
	Id            int
	Email         string `orm:"unique"`
	Password      string
	Registertime  time.Time   `orm:"auto_now_add;type(datetime)"`
	Lastlogintime time.Time   `orm:"type(datetime);null" form:"-"`
	IsActive      bool        `orm:"default(false)"`
	Profile       *Profile    `orm:"rel(one);unique"` // OneToOne relation
	UserData      []*UserData `orm:"reverse(many)"`   // reverse relationship of fk
}

type UserData struct {
	Id      int
	Type    string
	Data    string    `orm:"type(text)"`
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
	User    *User     `orm:"rel(fk)"` // RelForeignKey relation
}

// Profile 使用者個人資料
type Profile struct {
	Id                            int
	Name                          string
	Identity                      string
	Points                        float64 `orm:"digits(12);decimals(2);default(0.00)"`
	Skype                         string
	User                          *User                            `orm:"reverse(one);on_delete(set_null)"` // Reverse relationship (optional)
	Student                       *Student                         `orm:"null;rel(one)"`                    // Reverse relationship (optional)
	Teacher                       *Teacher                         `orm:"null;rel(one)"`                    // Reverse relationship (optional)
	Schedules                     []*LessonSchedule                `orm:"reverse(many)"`                    // reverse relationship of fk
	EZPayPaymentApplicationRecord []*EZPayPaymentApplicationRecord `orm:"reverse(many)"`                    // reverse relationship of fk
}

type Teacher struct {
	Id               int
	TeachingYears    int
	Proficiency      string `orm:"type(text)"`
	Resume           string
	IsActive         bool               `orm:"default(false)"`
	Profile          *Profile           `orm:"reverse(one)"`  // Reverse relationship (optional)
	TeacherTags      *TeacherTags       `orm:"rel(one)"`      // OneToOne relation
	StudentAuditings []*StudentAuditing `orm:"reverse(many)"` // reverse relationship of fk
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

type Student struct {
	Id               int
	AduitingTimes    int                `orm:"default(0)"`
	LeaveNumber      int                `orm:"default(0)"`
	Profile          *Profile           `orm:"reverse(one)"`  // Reverse relationship (optional)
	StudentAuditings []*StudentAuditing `orm:"reverse(many)"` // reverse relationship of fk
}

func init() {
	// Need to register model in init
	orm.RegisterModel(new(User), new(Profile), new(Teacher), new(TeacherTags), new(UserData), new(Student))

}
