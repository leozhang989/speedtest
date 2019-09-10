package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Users struct {
	Id                    int64  `orm:"auto"`
	DeviceCode            string `orm:"size(100)"`
	VipExpirationTime     uint64
	OriginalTransactionId string `orm:"size(30)"`
	Updated               uint64
	Created               uint64
}

func init() {
	orm.RegisterModel(new(Users))
}

// AddUsers insert a new Users into database and returns
// last inserted Id on success.
func AddUsers(m *Users) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUsersById retrieves Users by Id. Returns error if
// Id doesn't exist
func GetUsersById(id int64) (v *Users, err error) {
	o := orm.NewOrm()
	v = &Users{Id: id}
	if err = o.QueryTable(new(Users)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetUsersByDeviceCode retrieves Users by deviceCode. Returns error if
// Id doesn't exist
func GetUsersByDeviceCode(deviceCode string)(v *Users, err error)  {
	o := orm.NewOrm()
	v = &Users{DeviceCode: deviceCode}
	if err = o.QueryTable(new(Users)).Filter("device_code", deviceCode).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetUsersByOtid retrieves Users by deviceCode. Returns error if
// Id doesn't exist
func GetUsersByOtid(originalTransactionId string)(num int64, err error)  {
	o := orm.NewOrm()
	var users []*Users
	if num, err := o.QueryTable("users").Filter("original_transaction_id", originalTransactionId).All(&users); err == nil {
		return num, nil
	}

	return num, err
}

// GetAllUsers retrieves all Users matches certain condition. Returns empty list if
// no records exist
func GetAllUsers(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Users))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Users
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateUsersById updates Users by Id and returns error if
// the record to be updated doesn't exist
func UpdateUsersById(m *Users) (err error) {
	o := orm.NewOrm()
	v := Users{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// UpdateUsers updates Users by OriginalTransactionId and returns error if
// the record to be updated doesn't exist
func UpdateUsersByOtid(etime uint64, originalTransactionId string) (num int64, err error) {
	o := orm.NewOrm()
	if etime != 0 && len(originalTransactionId) != 0 {
		num, err := o.QueryTable("users").Filter("original_transaction_id", originalTransactionId).Update(orm.Params{
			"vip_expiration_time": etime, "updated": uint64(time.Now().Unix()),
		})
		fmt.Println("Number of records updated in database:", num)
		return num, err
	}

	return num, err
}

// UpdateUserInfoByOtid updates Users by OriginalTransactionId and returns error if
// the record to be updated doesn't exist
func UpdateUserInfoByOtid(m *Users) (num int64, err error) {
	o := orm.NewOrm()
	v := Users{OriginalTransactionId: m.OriginalTransactionId}
	// ascertain id exists in the database
	if err = o.Read(&v, "original_transaction_id"); err == nil {
		if num, err := o.QueryTable("users").Filter("original_transaction_id", m.OriginalTransactionId).Update(orm.Params{
			"vip_expiration_time": m.VipExpirationTime, "updated": m.Updated,
		}); err == nil{
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUsers deletes Users by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUsers(id int64) (err error) {
	o := orm.NewOrm()
	v := Users{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Users{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
