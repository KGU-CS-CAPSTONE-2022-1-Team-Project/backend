package dao

import (
	"backend/tool"
	"github.com/kamva/mgm/v3"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dsn string

func init() {
	if dsn == "" {
		tool.ReadConfig("./configs/partner", "secret", "json")
		dsnInfo := viper.GetStringMapString("db")
		dsn = "mongodb://" + dsnInfo["user"] + ":" + dsnInfo["password"] + "@" + dsnInfo["host"] + ":" + dsnInfo["port"]
		if err := mgm.SetDefaultConfig(nil, dsnInfo["name"], options.Client().ApplyURI(dsn)); err != nil {
			panic(err)
		}
	}
}
