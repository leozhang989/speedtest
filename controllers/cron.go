package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"net"
	"os"
)

// CronController operations for Cron
type CronController struct {
	beego.Controller
}


//Refresh method handles GET requests for CronController.
func (c *CronController) Refresh() {
	c.EnableRender = false
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println(ipnet.IP.String())
			}
		}
	}
}
