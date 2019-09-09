package controllers

import (
	"speedtest/models"
	"time"

	"github.com/astaxie/beego"
)

//  UsersController operations for Users
type UsersController struct {
	beego.Controller
}

type Result struct {
	Code int
	Data map[string]string
	Msg string
}

// URLMapping ...
func (c *UsersController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Users
// @Param	body		body 	models.Users	true		"body for Users content"
// @Success 201 {int} models.Users
// @Failure 403 body is empty
// @router /users [post]
func (c *UsersController) Post() {

}

// GetOne ...
// @Title Get One
// @Description get Users by deviceCode
// @Param	deviceCode		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Users
// @Failure 403 :deviceCode is empty
// @router /users/:deviceCode [get]
func (c *UsersController) GetOne() {
	c.EnableRender = false
	//idStr := c.Ctx.Input.Param(":id")
	//id, _ := strconv.ParseInt(idStr, 0, 64)
	deviceCodeStr := c.Ctx.Input.Param(":deviceCode")
	v, err := models.GetUsersByDeviceCode(deviceCodeStr)
	res := new(Result)
	settings, seterr := models.GetSettingsBySettingKey("downloadUrl")
	downloadUrl := ""
	if seterr == nil {
		downloadUrl = settings.SettingValue
	}
	now := time.Now().Unix()
	var vipetime string = "已过期"
	var isvip string = "0"
	if err != nil {
		//新用户创建
		var u models.Users
		u.DeviceCode = deviceCodeStr
		u.VipExpirationTime = uint64(now)
		u.OriginalTransactionId = ""
		u.Updated = uint64(now)
		u.Created = uint64(now)
		if _, err := models.AddUsers(&u);err != nil {
			//c.Data["json"] = err.Error()
			c.Ctx.Output.SetStatus(202)
			res.Code = 202
			res.Data = make(map[string]string)
			res.Msg = "login failed"
			c.Data["json"] = res
			c.ServeJSON()
		}
		vipetime = "已过期"
	} else {
		//用时间模板  格式化时间戳  时间转化为string，layout必须为 "2006-01-02 15:04:05"
		vipetime = time.Unix(int64(v.VipExpirationTime), 0).Format("2006-01-02 15:04:05")
		if now < int64(v.VipExpirationTime) {
			isvip = "1"
		}
	}
	returnRes := map[string]string{"VipExpirationTime": vipetime, "downloadUrl": downloadUrl, "IsVip": isvip}
	c.Ctx.Output.SetStatus(200)
	res.Code = 200
	res.Data = returnRes
	res.Msg = "success"
	c.Data["json"] = res
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Users
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Users
// @Failure 403
// @router /users [get]
func (c *UsersController) GetAll() {
	c.EnableRender = false
}

// Put ...
// @Title Put
// @Description update the Users
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Users	true		"body for Users content"
// @Success 200 {object} models.Users
// @Failure 403 :id is not int
// @router /users/:id [put]
func (c *UsersController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Users
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /users/:id [delete]
func (c *UsersController) Delete() {

}
