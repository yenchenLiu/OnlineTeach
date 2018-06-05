package models

import (
	"github.com/astaxie/beego/orm"
)

func (s *Student) Read(fields ...string) error {
	if err := orm.NewOrm().Read(s, fields...); err != nil {
		return err
	}
	return nil
}

func (s *Student) LoadAuditing() error {
	if _, err := orm.NewOrm().LoadRelated(s, "StudentAuditings"); err != nil {
		return err
	}
	return nil
}
