package models

import (
	"github.com/astaxie/beego/orm"
)

func (this *EZPayPaymentApplicationRecord) Read(fields ...string) error {
	if err := orm.NewOrm().Read(this, fields...); err != nil {
		return err
	}
	return nil
}

func (this *EZPayPaymentApplicationRecord) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(this, fields...); err != nil {
		return err
	}
	return nil
}

func (this *EZPayPaymentApplicationRecord) Insert() error {
	if _, err := orm.NewOrm().Insert(this); err != nil {
		return err
	}
	return nil
}