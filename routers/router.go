package routers

import (
	"github.com/astaxie/beego"
	"speedtest/controllers"
)

func init() {
    //beego.Router("/", &controllers.MainController{})
    //beego.Get("/user", func(ctx *context.Context) {
    //	ctx.Output.Body([]byte("hello world"))
	//})
	beego.Include(&controllers.UsersController{})
	beego.Include(&controllers.OrdersController{})
}
