package dao

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var config map[string]string

func init() {
	if !viper.IsSet("db") {
		viper.SetConfigName("client_secret")
		viper.SetConfigType("json")
		viper.AddConfigPath("./configs/owner")
		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("viper error: %v", err))
		}
	}
	config = viper.GetStringMapString("db")
}

func dbConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config["user"], config["password"], config["host"], config["port"], config["main_db"])
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db, err
}
