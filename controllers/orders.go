package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"speedtest/models"
	"strconv"
	"time"

	"github.com/astaxie/beego"
)

//  OrdersController operations for Orders
type OrdersController struct {
	beego.Controller
}

type OrdersResult struct {
	Code int
	Data map[string]string
	Msg string
}

//苹果内购返回内容解析-------
type In_app struct {
	Product_id string
}

type Receipt struct {
	Bundle_id string
	In_app []In_app
}

type Latest_receipt_info struct {
	Original_transaction_id string
	Expires_date_ms string
}

type AppleResult struct {
	Status int
	Environment string
	Latest_receipt_info []Latest_receipt_info
	Receipt Receipt
	Latest_receipt string
}
//苹果内购返回内容解析-------

// URLMapping ...
func (c *OrdersController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Orders
// @Param	body		body 	models.Orders	true		"body for Orders content"
// @Success 201 {int} models.Orders
// @Failure 403 body is empty
// @router /orders [post]
func (c *OrdersController) Post() {
	c.EnableRender = false
	var v models.Orders
	res := new(OrdersResult)
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if len(v.DeviceCode) == 0 || len(v.Certificate) == 0 {
		c.Ctx.Output.SetStatus(202)
		res.Code = 202
		res.Data = make(map[string]string)
		res.Msg = "wrong params"
		c.Data["json"] = res
		c.ServeJSON()
		panic("")
	}
	user,err := models.GetUsersByDeviceCode(v.DeviceCode)
	if err != nil || user == nil {
		c.Ctx.Output.SetStatus(202)
		res.Code = 202
		res.Data = make(map[string]string)
		res.Msg = "用户信息获取失败，请退出APP重试"
		c.Data["json"] = res
		c.ServeJSON()
		panic("")
	}
	isSandBoxRes, _ := models.GetSettingsBySettingKey("isSandBox")
	appPasswordRes, _ := models.GetSettingsBySettingKey("appPassword")
	bundleIdRes, _ := models.GetSettingsBySettingKey("bundleId")

	appleParams := `{"receipt-data":"` + v.Certificate + `", "password":"` + appPasswordRes.SettingValue + `"}`
	appleVerifyHost := "https://buy.itunes.apple.com/verifyReceipt"
	if intSandbox,_ := strconv.Atoi(isSandBoxRes.SettingValue); intSandbox == 1 {
		appleVerifyHost = "https://sandbox.itunes.apple.com/verifyReceipt";
	}

	appleResponseData,err := HttpPostJson(appleParams, appleVerifyHost)
	//fmt.Println(string(appleResponseData))
	if err != nil {
		c.Ctx.Output.SetStatus(202)
		res.Code = 202
		res.Data = make(map[string]string)
		res.Msg = "get result from apple error"
		c.Data["json"] = res
		c.ServeJSON()
		panic("")
	}
	var appleResponse AppleResult
	json.Unmarshal(appleResponseData, &appleResponse)
	//fmt.Println(appleResponse)

	switch {
	case appleResponse.Status == 21007:
		appleResponseData,err = HttpPostJson(appleParams, "https://sandbox.itunes.apple.com/verifyReceipt")

	case appleResponse.Status == 21008:
		appleResponseData,err = HttpPostJson(appleParams, "https://buy.itunes.apple.com/verifyReceipt")

	case appleResponse.Status >= 21100 && appleResponse.Status <= 21199:
		appleResponseData,err = HttpPostJson(appleParams, appleVerifyHost)
	}

	products := []string{"com.speed.1month", "com.speed.1year"}
	if appleResponse.Status == 0 {
		in_product_flag := false
		for _,v := range products{
			if v == appleResponse.Receipt.In_app[0].Product_id{
				in_product_flag = true
			}
		}
		if len(appleResponse.Receipt.Bundle_id) != 0 && appleResponse.Receipt.Bundle_id == bundleIdRes.SettingValue && len(appleResponse.Receipt.In_app) != 0 && in_product_flag {
			receiptInfoCount := len(appleResponse.Latest_receipt_info)
			if receiptInfoCount != 0 {
				lastestOrder := appleResponse.Latest_receipt_info[receiptInfoCount - 1]
				originalTransactionId := lastestOrder.Original_transaction_id
				userInfos, _ := models.GetUsersByOtid(originalTransactionId)
				now := time.Now().Unix()
				nowMicroS := now * 1000
				appleVipExpires, _ := strconv.ParseInt(lastestOrder.Expires_date_ms, 10, 64)
				if nowMicroS > appleVipExpires {
					c.Ctx.Output.SetStatus(202)
					res.Code = 202
					res.Data = make(map[string]string)
					res.Msg = "续订已过期，请重新购买"
					c.Data["json"] = res
					c.ServeJSON()
					panic("")
				}

				//获取orders表中是否有此设备记录
				orderRecord, _ := models.GetOrdersByDeviceCode(v.DeviceCode)
				o := orm.NewOrm()
				ormerr := o.Begin()
				var doerrs error
				if orderRecord == nil {
					//不存在时则创建新的原始订单记录 存储在 orders 表中
					v.Updated = uint64(now)
					v.PayStatus = true
					v.LatestReceipt = appleResponse.Latest_receipt
					v.Created = uint64(now)
					//续订成功时，会员原剩余时长保存，续订结束时继续使用
					_, doerrs = models.AddOrders(&v)
				}else{
					orderRecord.Updated = uint64(now)
					orderRecord.PayStatus = true
					orderRecord.LatestReceipt = appleResponse.Latest_receipt
					doerrs = models.UpdateOrdersById(orderRecord)
				}

				//更新会员到期时间
				expiresDateS, _ := strconv.ParseUint(lastestOrder.Expires_date_ms[:len(lastestOrder.Expires_date_ms)-3], 10, 64)
				user.VipExpirationTime = expiresDateS
				user.OriginalTransactionId = originalTransactionId
				user.Updated = uint64(now)
				doerrs = models.UpdateUsersById(user)
				if userInfos != 0 {
					_, doerrs = models.UpdateUsersByOtid(expiresDateS, originalTransactionId)
				}

				if doerrs != nil || ormerr != nil {
					ormerr = o.Rollback()
				}else {
					ormerr = o.Commit()
					c.Ctx.Output.SetStatus(200)
					res.Code = 200
					vipetime := time.Unix(int64(expiresDateS), 0).Format("2006-01-02 15:04:05")
					dataRes := map[string]string{"IsVip":"1", "VipExpirationTime":vipetime}
					res.Data = dataRes
					res.Msg = "success"
					c.Data["json"] = res
					c.ServeJSON()
				}
			}
		}
	}

	c.Ctx.Output.SetStatus(202)
	res.Code = 202
	res.Data = make(map[string]string)
	res.Msg = "支付失败"
	c.Data["json"] = res
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Orders by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Orders
// @Failure 403 :id is empty
// @router /orders/:id [get]
func (c *OrdersController) GetOne() {

}

// GetAll ...
// @Title Get All
// @Description get Orders
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Orders
// @Failure 403
// @router /orders [get]
func (c *OrdersController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Orders
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Orders	true		"body for Orders content"
// @Success 200 {object} models.Orders
// @Failure 403 :id is not int
// @router /orders/:id [put]
func (c *OrdersController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Orders
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /orders/:id [delete]
func (c *OrdersController) Delete() {

}
