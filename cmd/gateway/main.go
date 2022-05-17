package main

import (
	"backend/api/gateway/owner"
	"backend/tool"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	Auth    = "/auth"
	Address = "/owner/address"
)

var port string

func init() {
	if !viper.IsSet("port") {
		tool.ReadConfig("./configs/gateway", "services", "yaml")
	}
	port = viper.GetStringMapString("gateway")["port"]
}

func main() {
	r := gin.Default()
	r.GET(Auth, owner.IsAuth)
	r.PUT(Address, owner.Address)
	err := r.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
