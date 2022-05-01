package dao

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func dbConnection() (*gorm.DB, error) {
	config := viper.GetStringMapString("db")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config["user"], config["password"], config["host"], config["port"], config["main_db"])
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db, err
}
