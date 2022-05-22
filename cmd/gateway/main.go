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
	Nickname           = "/common/nickname"
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
	r.POST(Nickname, owner.SetNickname)
	r.GET(Nickname+"/:address", owner.GetNickname)
	err := r.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
