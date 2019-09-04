package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "speedtest/routers"
)

func main() {
	dbuser := beego.AppConfig.String("mysqluser")
	dbpass := beego.AppConfig.String("mysqlpass")
	dbmysqlurls := beego.AppConfig.String("mysqlurls")
	dbmysqldb := beego.AppConfig.String("mysqldb")
	dbInfo := dbuser + ":" + dbpass + "@tcp(" + dbmysqlurls + ":3306)/" + dbmysqldb
	orm.RegisterDataBase("default","mysql",dbInfo)

	beego.Run()
}

