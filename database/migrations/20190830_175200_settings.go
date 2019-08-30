package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Settings_20190830_175200 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Settings_20190830_175200{}
	m.Created = "20190830_175200"

	migration.Register("Settings_20190830_175200", m)
}

// Run the migrations
func (m *Settings_20190830_175200) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE settings(`id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',`setting_key` varchar(255) NOT NULL DEFAULT '',`setting_value` varchar(255) NOT NULL DEFAULT '',PRIMARY KEY (`id`))")
}

// Reverse the migrations
func (m *Settings_20190830_175200) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `settings`")
}
