package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Users_20190830_183829 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Users_20190830_183829{}
	m.Created = "20190830_183829"

	migration.Register("Users_20190830_183829", m)
}

// Run the migrations
func (m *Users_20190830_183829) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE users(`id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',`device_code` varchar(100) NOT NULL DEFAULT '',`vip_expiration_time` int(11) NOT NULL DEFAULT 0,`original_transaction_id` varchar(30) NOT NULL DEFAULT '',`updated` int(11) NOT NULL DEFAULT 0,`created` int(11) NOT NULL DEFAULT 0,PRIMARY KEY (`id`))")
}

// Reverse the migrations
func (m *Users_20190830_183829) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `users`")
}
