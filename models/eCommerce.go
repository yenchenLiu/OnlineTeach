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
	Created           time.Time                    `orm:"auto_now_add;type(datetime)"`
	Updated           time.Time                    `orm:"auto_now;type(datetime)"`
	Profile           *Profile                     `orm:"rel(fk);null"`        // RelForeignKey relation
	ReceiveRecord     []*EZPayPaymentReceiveRecord `orm:"reverse(many); null"` // reverse relationship of fk

}

type EZPayPaymentReceiveRecord struct {
	Id                       int
	MerchantID               string
	MerchantTradeNo          string
	RtnCode                  string
	RtnMsg                   string
	TradeNo                  string
	TradeAmt                 string
	PaymentDate              string
	PaymentType              string
	PaymentTypeChargeFee     string
	TradeDate                string
	SimulatePaid             string
	PaymentApplicationRecord *EZPayPaymentApplicationRecord `orm:"rel(fk)"`       // RelForeignKey relation
	PointsTrade              *PointsTrade                   `orm:"rel(one);null"` // OneToOne relation
}

type PointsTrade struct {
	Id                   int
	Points               float64                    `orm:"digits(12);decimals(2);default(0.00)"`
	Description          string                     `orm:"type(text)"`
	PaymentReceiveRecord *EZPayPaymentReceiveRecord `orm:"reverse(one);null"` // Reverse relationship (optional)
	CourseRecord         *CourseRecord              `orm:"rel(fk);null"`      // Reverse relationship (optional)
	ProfileReceiver      *Profile                   `orm:"rel(fk);null"`      // RelForeignKey relation
	ProfileGiver         *Profile                   `orm:"rel(fk);null"`      // RelForeignKey relation
	Created              time.Time                  `orm:"auto_now_add;type(datetime)"`
	Updated              time.Time                  `orm:"auto_now;type(datetime)"`
}

type WithdrawRecord struct {
	Id          int
	Points      float64 `orm:"digits(12);decimals(2);default(0.00)"`
	PayPal      string
	Profile     *Profile  `orm:"rel(fk);null"` // RelForeignKey relation
	Process    string    // 處理進度 [suspending、handling、solved、error]
	Description string    `orm:"type(text)"`
	Created     time.Time `orm:"auto_now_add;type(datetime)"`
	Updated     time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	// Need to register model in init
	orm.RegisterModel(new(EZPayPaymentApplicationRecord), new(EZPayPaymentReceiveRecord), new(PointsTrade), new(WithdrawRecord))

}
