package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"speedtest/controllers"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Get("/user", func(ctx *context.Context) {
    	ctx.Output.Body([]byte("hello world"))
	})
}
