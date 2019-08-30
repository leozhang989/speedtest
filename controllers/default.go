package controllers

import (
	"github.com/astaxie/beego"
	"strings"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	kv := strings.SplitN("string:255:''", ":", 3)
	var flag = 0
	if len(kv) == 3 {
		flag = 1
	}
	//c.Data["Website"] = "beego.me"
	c.Data["Website"] = flag
	c.Data["Email"] = "powerof1024@gmail.com"
	c.TplName = "index.tpl"
}
