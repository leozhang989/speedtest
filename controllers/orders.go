package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"net/http"
	"speedtest/models"
	"strconv"
	"strings"
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
		res.Msg = "empty user"
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

	appleResponseData,err := httpPostJson(appleParams, appleVerifyHost)
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
		appleResponseData,err = httpPostJson(appleParams, "https://sandbox.itunes.apple.com/verifyReceipt")

	case appleResponse.Status == 21008:
		appleResponseData,err = httpPostJson(appleParams, "https://buy.itunes.apple.com/verifyReceipt")

	case appleResponse.Status >= 21100 && appleResponse.Status <=21199:
		appleResponseData,err = httpPostJson(appleParams, appleVerifyHost)
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
				userInfo,_ := models.GetUsersByOtid(originalTransactionId)
				if userInfo == nil {
					now := time.Now().Unix()
					if string(now * 1000) > lastestOrder.Expires_date_ms {
						c.Ctx.Output.SetStatus(202)
						res.Code = 202
						res.Data = make(map[string]string)
						res.Msg = "续订已过期，请重新购买"
						c.Data["json"] = res
						c.ServeJSON()
						panic("")
					}
					//不存在时则创建新的原始订单记录 存储在 orders 表中
					o := orm.NewOrm()
					ormerr := o.Begin()
					v.Created = uint64(now)
					v.Updated = uint64(now)
					v.PayStatus = true
					//续订成功时，会员原剩余时长保存，续订结束时继续使用
					_, doerrs := models.AddOrders(&v)

					//更新会员到期时间
					expiresDateS, _ := strconv.ParseUint(lastestOrder.Expires_date_ms[:len(lastestOrder.Expires_date_ms)-3], 10, 64)
					user.VipExpirationTime = expiresDateS
					user.OriginalTransactionId = originalTransactionId
					doerrs = models.UpdateUsersById(user)

					if doerrs != nil || ormerr != nil {
						ormerr = o.Rollback()
					}else {
						ormerr = o.Commit()
					}
				}else {
					expiresDateS, _ := strconv.ParseUint(lastestOrder.Expires_date_ms[:len(lastestOrder.Expires_date_ms)-3], 10, 64)
					user.VipExpirationTime = expiresDateS
					models.UpdateUsersById(user)
				}
				//更新最新凭证
				v.LatestReceipt = appleResponse.Latest_receipt
				models.UpdateOrdersById(&v)
				c.Ctx.Output.SetStatus(200)
				res.Code = 200
				res.Data = make(map[string]string)
				res.Msg = "success"
				c.Data["json"] = res
				c.ServeJSON()
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
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v, err := models.GetOrdersById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
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
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllOrders(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
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
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v := models.Orders{Id: id}
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err := models.UpdateOrdersById(&v); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Orders
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /orders/:id [delete]
func (c *OrdersController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := models.DeleteOrders(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func httpPostJson(requestJsonString string, requestUrl string) ([]byte, error) {
	jsonByteStr :=[]byte(requestJsonString)
	req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(jsonByteStr))
	if err != nil {
		return nil,errors.New("request error")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil,errors.New("request error")
	}
	defer resp.Body.Close()

	//statuscode := resp.StatusCode
	//header := resp.Header
	body,_ := ioutil.ReadAll(resp.Body)
	//response := map[string]string{"code":strconv.Itoa(statuscode),"body":string(body[:])}
	return body,nil
}