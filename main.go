package main

import (
	"WebPartice/models"
	_ "WebPartice/routers"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	orm.RegisterDriver("sqlite3", orm.DRSqlite)
	orm.RegisterDataBase("default", "sqlite3", "./test.db")
	fmt.Println("test")
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
	profile.Name = "yenchen"

	user := new(models.User)
	user.Profile = profile
	user.Email = "mail@daychen.tw"
	user.Password = "yenchen"

	fmt.Println(o.Insert(profile))
	fmt.Println(o.Insert(user))
	beego.Run()
}
