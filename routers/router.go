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
	//注解路由，详见controller中的注释
	beego.Include(&controllers.UsersController{})
	beego.Include(&controllers.OrdersController{})

	// cron定时任务
	beego.Router("/cron/refresh-vip/:token", &controllers.CronController{}, "get:Refresh")
}
