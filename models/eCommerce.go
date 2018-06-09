package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type EZPayPaymentApplicationRecord struct {
	Id                int
	MerchantTradeNo   string `orm:"unique"`
	MerchantTradeDate string
	TotalAmount       int
	ReturnURL         string
	Created           time.Time `orm:"auto_now_add;type(datetime)"`
	Updated           time.Time `orm:"auto_now;type(datetime)"`
	Profile           *Profile  `orm:"rel(fk);null"` // RelForeignKey relation

}

func init() {
	// Need to register model in init
	orm.RegisterModel(new(EZPayPaymentApplicationRecord))

}
