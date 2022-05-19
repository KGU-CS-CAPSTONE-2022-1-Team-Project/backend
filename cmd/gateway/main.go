package main

import (
	"backend/api/gateway/owner"
	"backend/api/gateway/partner"
	"backend/tool"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	GoogleAuth         = "/owner/google"
	GoogleAuthCallback = "/owner/google/callback"
	Youtuber           = "/owner/youtuber"
	Nft                = "/partner/nft"
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
	r.GET(GoogleAuth, owner.Redirecting)
	r.GET(GoogleAuthCallback, owner.RegisterUser)
	r.PUT(Youtuber, owner.AuthYoutuber)
	r.GET(Youtuber+"/:id", owner.GetChannel)
	r.POST(Nft, partner.CheckFile, partner.Upload)
	r.GET(Nft+"/:id", partner.GetNFTInfo)
	err := r.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
