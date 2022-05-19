package dao

import (
	"backend/tool"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dsn string

func init() {
	if dsn == "" {
		tool.ReadConfig("configs/owner", "client_secret", "json")
		config := viper.GetStringMapString("db")
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config["user"], config["password"], config["host"], config["port"], config["main_db"])
	}
}

func dbConnection() (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
