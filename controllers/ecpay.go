package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/astaxie/beego"
)

type ECPayController struct {
	BaseController
}

func (this *ECPayController) Prepare() {
	this.EnableXSRF = false
}

func (this *ECPayController) Post() {
	hashKey := beego.AppConfig.String("ECPAYHashKey")
	hashIV := beego.AppConfig.String("ECPAYHashIV")
	fmt.Println("Receive Data", this.Input())

	m := make(map[string]string)
	m["MerchantID"] = this.Input()["MerchantID"][0]
	m["MerchantTradeNo"] = this.Input()["MerchantTradeNo"][0]
	m["StoreID"] = ""
	m["RtnCode"] = this.Input()["RtnCode"][0] //Check this number is 1, it mean is success
	m["RtnMsg"] = this.Input()["RtnMsg"][0]
	m["TradeNo"] = this.Input()["TradeNo"][0]
	m["TradeAmt"] = this.Input()["TradeAmt"][0]
	m["PaymentDate"] = this.Input()["PaymentDate"][0]
	m["PaymentType"] = this.Input()["PaymentType"][0]
	m["PaymentTypeChargeFee"] = this.Input()["PaymentTypeChargeFee"][0]
	m["TradeDate"] = this.Input()["TradeDate"][0]
	m["SimulatePaid"] = this.Input()["SimulatePaid"][0] //Check this number is 0, it mean is true trade
	m["CustomField1"] = ""
	m["CustomField2"] = ""
	m["CustomField3"] = ""
	m["CustomField4"] = ""
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	t := ""
	for _, k := range keys {
		t += k + "=" + m[k] + "&"
	}
	CheckMacValue := "HashKey=" + hashKey + "&"
	CheckMacValue += t
	CheckMacValue += "HashIV=" + hashIV
	CheckMacValue = url.QueryEscape(CheckMacValue)
	CheckMacValue = strings.ToLower(CheckMacValue)
	h := sha256.New()
	h.Write([]byte(CheckMacValue))
	CheckMacValue = hex.EncodeToString(h.Sum(nil))
	CheckMacValue = strings.ToUpper(CheckMacValue)
	fmt.Println(this.Input()["CheckMacValue"][0])
	fmt.Println(CheckMacValue)
	fmt.Println(this.Input()["CheckMacValue"][0] == CheckMacValue)

	this.Data["response"] = "1|OK"
	this.TplName = "student/ecPayResponse.tpl"
}
