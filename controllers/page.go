package controllers

import (
	"github.com/astaxie/beego"
)

// PageController operations for Page
type PageController struct {
	beego.Controller
}


func (c *PageController) Service() {
	c.TplName = "page/service.tpl"
}

func (c *PageController) Policy() {
	c.TplName = "page/privacy-policy.tpl"
}