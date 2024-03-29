package main

import (
	"OnlineTeach/controllers"
	"OnlineTeach/lib"
	"OnlineTeach/models"
	_ "OnlineTeach/routers"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	orm.RegisterDriver("sqlite3", orm.DRSqlite)
	orm.RegisterDataBase("default", "sqlite3", "./test.db")
	force := false
	verbose := true
	err := orm.RunSyncdb("default", force, verbose)
	if err != nil {
		fmt.Println(err)
	}
}

func createSuperUser() {
	o := orm.NewOrm()
	o.Using("default")

	err := o.Begin()

	profile := new(models.Profile)
	profile.Name = "Yenchen"
	profile.Identity = "admin"

	user := new(models.User)
	user.Profile = profile
	user.Email = "mail@daychen.tw"
	user.Password = "yenchen"
	user.IsActive = true
	h := sha256.New()
	h.Write([]byte(user.Email))
	h.Write([]byte(user.Password))
	user.Password = string(base64.URLEncoding.EncodeToString(h.Sum(nil)))
	if err == nil {
		_, err = o.Insert(profile)
	}
	if err == nil {
		_, err = o.Insert(user)
	}
	if err != nil {
		o.Rollback()
	} else {
		o.Commit()
	}
}

func createCustodyAccount() {

	o := orm.NewOrm()
	o.Using("default")

	err := o.Begin()
	// 上課點數保管
	profile := new(models.Profile)
	profile.Name = "Custody"
	profile.Identity = "admin"

	user := new(models.User)
	user.Profile = profile
	user.Email = "custody@daychen.tw"
	user.Password = "yenchen"
	user.IsActive = false

	if err == nil {
		_, err = o.Insert(profile)
	}
	if err == nil {
		_, err = o.Insert(user)
	}
	if err != nil {
		o.Rollback()
	} else {
		o.Commit()
	}
}

func createWithdrawAccount() {

	o := orm.NewOrm()
	o.Using("default")

	err := o.Begin()
	// 點數提領帳戶
	profile := new(models.Profile)
	profile.Name = "Withdraw"
	profile.Identity = "admin"

	user := new(models.User)
	user.Profile = profile
	user.Email = "withdraw@daychen.tw"
	user.Password = "yenchen"
	user.IsActive = false

	if err == nil {
		_, err = o.Insert(profile)
	}
	if err == nil {
		_, err = o.Insert(user)
	}
	if err != nil {
		o.Rollback()
	} else {
		o.Commit()
	}
}

func main() {

	orm.Debug = true

	createSuperUser()
	createCustodyAccount()
	createWithdrawAccount()

	beego.AddFuncMap("AddNumber", lib.AddNumber)
	beego.ErrorController(&controllers.ErrorController{})
	beego.Run()
}
