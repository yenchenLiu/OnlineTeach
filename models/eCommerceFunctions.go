package models

import (
	"strconv"
	"github.com/astaxie/beego/orm"
)

func (this *EZPayPaymentApplicationRecord) Read(fields ...string) error {
	if err := orm.NewOrm().Read(this, fields...); err != nil {
		return err
	}
	return nil
}

func (this *EZPayPaymentApplicationRecord) LoadProfile() error {
	if _, err := orm.NewOrm().LoadRelated(this, "Profile"); err != nil {
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

func (this *EZPayPaymentReceiveRecord) Insert() error {
	if _, err := orm.NewOrm().Insert(this); err != nil {
		return err
	}
	return nil
}

func (this *EZPayPaymentReceiveRecord) Deposit(p *Profile) error {
	money, _ := strconv.ParseInt(this.TradeAmt,10,64)
	points := money / 100
	o := orm.NewOrm()
	err := o.Begin()
	trade := new(PointsTrade)
	trade.Points = float64(points)
	trade.Description = "儲值" + this.TradeAmt + "元"
	trade.Profile = p
	if _, err = o.Insert(trade); err != nil {
		o.Rollback()
		return err
	}

	p.Points += float64(points)
	if _, err = o.Update(p, "Points"); err != nil {
		o.Rollback()
		return err
	}
	this.PointsTrade = trade
	if _, err = o.Insert(this); err != nil {
		o.Rollback()
		return err
	}

	o.Commit()
	return nil
}