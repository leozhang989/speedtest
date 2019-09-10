package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
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

	orders, err := models.GetAllOrders(queryCondition, queryFields, sortBy, orderBy, 0, 0, nil)
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

		//如果已经被前面的更新过则不再更新 检查更新时间
		checkStatus := CheckUpdatedTime(now, order)
		if !checkStatus {
			continue
		}

		res, err := DealAppleOrder(order.LatestReceipt)
		//logs.Info("【定时任务】苹果接口返回错误记录：订单id是%v，错误是%s，返回的结果是%s", order.Id, err, res)
		if err != nil || len(res) == 0 {
			//记录错误日志
			logs.Info("【定时任务】苹果接口返回错误记录：订单id是%v，错误是%s，返回的结果是%s", order.Id, err, res)
			//fmt.Println(err)
		}

		o := orm.NewOrm()
		ormerr := o.Begin()
		var doerrs error
		//更新会员到期时间和LatestReceipt 通过OriginalTransactionId批量更新
		var user models.Users
		expiresDateS, _ :=strconv.ParseUint(res["expiresDateS"], 10, 64)
		user.OriginalTransactionId = res["originalTransactionId"]
		user.VipExpirationTime = expiresDateS
		user.Updated = uint64(now)
		_, doerrs = models.UpdateUserInfoByOtid(&user)


		//更新订单表的最近凭证
		order.LatestReceipt = res["latestReceipt"]
		order.Updated = uint64(now)
		orderErr := models.UpdateOrdersById(&order)

		if doerrs != nil || ormerr != nil || orderErr != nil {
			ormerr = o.Rollback()
			fmt.Println("error")
		}else {
			ormerr = o.Commit()
			fmt.Println("success")
		}
	}
}
