package controllers

import (
	"OnlineTeach/models"
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
	this.TplName = "student/ecPayResponse.tpl"
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

	// Checksum錯誤，記錄至資料庫
	if this.Input()["CheckMacValue"][0] != CheckMacValue {
		record := &models.EZPayPaymentApplicationRecord{MerchantTradeNo: m["MerchantTradeNo"]}
		record.Read("MerchantTradeNo")
		receive := new(models.EZPayPaymentReceiveRecord)
		receive.MerchantID = m["MerchantID"]
		receive.MerchantTradeNo = m["MerchantTradeNo"]
		receive.RtnCode = m["RtnCode"]
		receive.RtnMsg = m["RtnMsg"]
		receive.TradeNo = m["TradeNo"]
		receive.TradeAmt = m["TradeAmt"]
		receive.PaymentDate = m["PaymentDate"]
		receive.PaymentType = m["PaymentType"]
		receive.PaymentTypeChargeFee = m["PaymentTypeChargeFee"]
		receive.TradeDate = m["TradeDate"]
		receive.SimulatePaid = m["SimulatePaid"]
		receive.PaymentApplicationRecord = record
		receive.Insert()
		return
	}

	// 如果交易失敗，或是在非開發模式下卻使用模擬付款則只記錄資料
	if (m["RtnCode"] != "1") || (m["SimulatePaid"] != "0" && beego.AppConfig.String("runmode") != "dev") {
		record := &models.EZPayPaymentApplicationRecord{MerchantTradeNo: m["MerchantTradeNo"]}
		record.Read("MerchantTradeNo")
		receive := new(models.EZPayPaymentReceiveRecord)
		receive.MerchantID = m["MerchantID"]
		receive.MerchantTradeNo = m["MerchantTradeNo"]
		receive.RtnCode = m["RtnCode"]
		receive.RtnMsg = m["RtnMsg"]
		receive.TradeNo = m["TradeNo"]
		receive.TradeAmt = m["TradeAmt"]
		receive.PaymentDate = m["PaymentDate"]
		receive.PaymentType = m["PaymentType"]
		receive.PaymentTypeChargeFee = m["PaymentTypeChargeFee"]
		receive.TradeDate = m["TradeDate"]
		receive.SimulatePaid = m["SimulatePaid"]
		receive.PaymentApplicationRecord = record
		receive.Insert()
		this.Data["response"] = "1|OK"
		return
	}

	// 交易成功，儲值並記錄至資料庫
	record := &models.EZPayPaymentApplicationRecord{MerchantTradeNo: m["MerchantTradeNo"]}
	record.Read("MerchantTradeNo")
	receive := new(models.EZPayPaymentReceiveRecord)
	receive.MerchantID = m["MerchantID"]
	receive.MerchantTradeNo = m["MerchantTradeNo"]
	receive.RtnCode = m["RtnCode"]
	receive.RtnMsg = m["RtnMsg"]
	receive.TradeNo = m["TradeNo"]
	receive.TradeAmt = m["TradeAmt"]
	receive.PaymentDate = m["PaymentDate"]
	receive.PaymentType = m["PaymentType"]
	receive.PaymentTypeChargeFee = m["PaymentTypeChargeFee"]
	receive.TradeDate = m["TradeDate"]
	receive.SimulatePaid = m["SimulatePaid"]
	receive.PaymentApplicationRecord = record
	receive.Deposit(record.Profile)

	this.Data["response"] = "1|OK"

}
