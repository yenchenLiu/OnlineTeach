package main

import (
	"WebPartice/controllers"
	"WebPartice/models"
	_ "WebPartice/routers"
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
	// Error.
	force := true
	verbose := true
	err := orm.RunSyncdb("default", force, verbose)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	orm.Debug = true
	o := orm.NewOrm()
	o.Using("default")
	profile := new(models.Profile)
	profile.Name = "Yenchen"
	profile.Identity = "admin"

	user := new(models.User)
	user.Profile = profile
	user.Email = "mail@daychen.tw"
	user.Password = "yenchen"
	h := sha256.New()
	h.Write([]byte(user.Email))
	h.Write([]byte(user.Password))
	user.Password = string(base64.URLEncoding.EncodeToString(h.Sum(nil)))

	fmt.Println(o.Insert(profile))
	fmt.Println(o.Insert(user))

	beego.ErrorController(&controllers.ErrorController{})
	beego.Run()
}
