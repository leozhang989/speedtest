package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Orders_20190830_184415 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Orders_20190830_184415{}
	m.Created = "20190830_184415"

	migration.Register("Orders_20190830_184415", m)
}

// Run the migrations
func (m *Orders_20190830_184415) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE orders(`id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',`device_code` varchar(100) NOT NULL DEFAULT '',`pay_status` tinyint(1) NOT NULL,`certificate` longtext NOT NULL,`latest_receipt` longtext NOT NULL,`updated` int(11) NOT NULL DEFAULT 0,`created` int(11) NOT NULL DEFAULT 0,PRIMARY KEY (`id`))")
}

// Reverse the migrations
func (m *Orders_20190830_184415) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `orders`")
}
