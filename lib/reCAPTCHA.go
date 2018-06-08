package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/astaxie/beego"
)

func ReCAPTCHAVerify(response string) bool {
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
		url.Values{"secret": {beego.AppConfig.String("reCAPTCHAKey")}, "response": {response}})

	if err != nil {
		fmt.Println(err)
	}
	var v map[string]interface{}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		fmt.Println(err)
	}
	if v["success"] == "true" {
		return true
	} else {
		return false
	}

}
