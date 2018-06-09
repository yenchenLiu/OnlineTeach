package lib

import (
	"OnlineTeach/models"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

// 綠界付款
func PayMoney(Id int, money int) string {
	hashKey := beego.AppConfig.String("EZPAYHashKey")
	hashIV := beego.AppConfig.String("EZPAYHashIV")
	m := make(map[string]string)
	m["MerchantID"] = "3076959"
	applyTime := time.Now()
	m["MerchantTradeNo"] = strconv.Itoa(Id) + applyTime.Format("20060102150405")
	m["MerchantTradeDate"] = applyTime.Format("2006/01/02 15:04:05")
	m["PaymentType"] = "aio"
	m["TotalAmount"] = strconv.Itoa(money)
	m["TradeDesc"] = "線上英文家教平台"
	m["ItemName"] = "線上英文家教平台" + strconv.Itoa(money/100) + "點"
	m["ReturnURL"] = "https://www.onlineteach.asia/ecpay/receive"
	m["ChoosePayment"] = "ALL"
	m["EncryptType"] = "1"

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	t := ""
	for _, k := range keys {
		t += k + "=" + m[k] + "&"
	}
	post := t
	CheckMacValue := "HashKey=" + hashKey + "&"
	CheckMacValue += t
	CheckMacValue += "HashIV=" + hashIV
	fmt.Println(CheckMacValue)
	CheckMacValue = url.QueryEscape(CheckMacValue)
	CheckMacValue = strings.ToLower(CheckMacValue)
	h := sha256.New()
	h.Write([]byte(CheckMacValue))
	CheckMacValue = hex.EncodeToString(h.Sum(nil))
	CheckMacValue = strings.ToUpper(CheckMacValue)
	post += "CheckMacValue=" + CheckMacValue

	profile := &models.Profile{Id: Id}
	profile.Read("Id")

	applicationRecord := new(models.EZPayPaymentApplicationRecord)
	applicationRecord.MerchantTradeNo = m["MerchantTradeNo"]
	applicationRecord.MerchantTradeDate = m["MerchantTradeDate"]
	applicationRecord.TotalAmount = money
	applicationRecord.ReturnURL = m["MerchantTradeNo"]
	applicationRecord.Profile = profile

	applicationRecord.Insert()

	return post
}

func PayMoneyTest(Id int, money int) string {
	hashKey := "5294y06JbISpM5x9"
	hashIV := "v77hoKGq4kWxNNIS"
	m := make(map[string]string)
	m["MerchantID"] = "2000214"
	m["MerchantTradeNo"] = strconv.Itoa(Id) + time.Now().Format("20060102150405")
	m["MerchantTradeDate"] = time.Now().Format("2006/01/02 15:04:05")
	m["PaymentType"] = "aio"
	m["TotalAmount"] = strconv.Itoa(money)
	m["TradeDesc"] = "線上英文家教平台"
	m["ItemName"] = "線上英文家教平台"
	m["ReturnURL"] = "https://www.onlineteach.asia/ecpay/receive/"
	m["ChoosePayment"] = "ALL"
	m["EncryptType"] = "1"

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	t := ""
	for _, k := range keys {
		t += k + "=" + m[k] + "&"
	}
	post := t
	CheckMacValue := "HashKey=" + hashKey + "&"
	CheckMacValue += t
	CheckMacValue += "HashIV=" + hashIV
	fmt.Println(CheckMacValue)
	CheckMacValue = url.QueryEscape(CheckMacValue)
	CheckMacValue = strings.ToLower(CheckMacValue)
	h := sha256.New()
	h.Write([]byte(CheckMacValue))
	CheckMacValue = hex.EncodeToString(h.Sum(nil))
	CheckMacValue = strings.ToUpper(CheckMacValue)
	post += "CheckMacValue=" + CheckMacValue

	return post
}
