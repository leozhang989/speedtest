package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"speedtest/models"
	"strconv"
	"time"
)

// CronController operations for Cron
type CronController struct {
	beego.Controller
}


//Refresh method handles GET requests for CronController.
func (c *CronController) Refresh() {
	//不渲染页面
	c.EnableRender = false

	cronTokenUrl := c.Ctx.Input.Param(":token")
	cronTokenConf := beego.AppConfig.String("crontoken")
	if len(cronTokenUrl) == 0 || cronTokenUrl != cronTokenConf {
		fmt.Println("token error")
		panic("")
	}

	//获取符合条件的支付凭证
	now := time.Now().Unix()
	timeConditon := now - 3600 * 12
	queryCondition := map[string]string{"pay_status":"1", "updated__gte":strconv.FormatInt(timeConditon,10)}
	queryFields := []string{}
	orderBy := []string{"asc"}
	sortBy := []string{"updated"}

	orders, err := models.GetAllOrders(queryCondition, queryFields, sortBy, orderBy, 0, 0)
	if err != nil {
		fmt.Println(err)
		panic("")
	}
	//fmt.Println(orders)

	var order models.Orders

	for _, orderOri := range orders {
		//a := reflect.TypeOf(order) 查看数据类型 然后断言之后进行赋值
		//接口取出的数据首先进行类型断言
		order = orderOri.(models.Orders)
		//如果已经被前面的更新过则不再更新
		//if order.Updated >= uint64(now) {
		//	continue
		//}

		res, err := DealAppleOrder(order.LatestReceipt)
		if err != nil || len(res) == 0 {
			//记录错误日志
			fmt.Println(err)
		}

		o := orm.NewOrm()
		ormerr := o.Begin()
		var doerrs error
		//更新会员到期时间和LatestReceipt
		var user models.Users
		expiresDateS, _ :=strconv.ParseUint(res["expiresDateS"], 10, 64)
		user.Id = order.Id
		user.VipExpirationTime = expiresDateS
		user.Updated = uint64(now)
		doerrs = models.UpdateUsersById(&user)

		//更新订单表的最近凭证
		var orderm models.Orders
		orderm.LatestReceipt = res["latestReceipt"]

		if doerrs != nil || ormerr != nil {
			ormerr = o.Rollback()
		}else {
			ormerr = o.Commit()
			fmt.Println("success")
		}

	}


}
