package main

import (
	"backend/api/gateway/owner"
	"backend/api/gateway/partner"
	"backend/tool"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

const (
	GoogleAuth         = "/owner/google"
	GoogleAuthCallback = "/owner/google/callback"
	Youtuber           = "/owner/youtuber"
	Nft                = "/partner/nft"
	Nickname           = "/common/nickname"
)

var port string

var corsList []string

const (
	fullchain  = "/etc/letsencrypt/live/capston-blockapp.greenflamingo.dev/fullchain.pem"
	privatekey = "/etc/letsencrypt/live/capston-blockapp.greenflamingo.dev/privkey.pem"
)

func init() {
	if port == "" {
		tool.ReadConfig("./configs/gateway", "services", "yaml")
		port = viper.GetStringMapString("gateway")["port"]
	}
	if corsList == nil {
		tool.ReadConfig("./configs/gateway", "client_secret", "json")
		corsList = viper.GetStringSlice("cors_list")
	}
}

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: corsList,
		AllowMethods: []string{http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodOptions},
		AllowHeaders: []string{"Origin", "content-type", "Authorization"},
	}))
	r.GET(GoogleAuth, owner.Redirecting)
	r.GET(GoogleAuthCallback, owner.RegisterUser)
	r.PUT(Youtuber, owner.AuthYoutuber)
	r.GET(Youtuber+"/:id", owner.GetChannel)
	r.POST(Nft, partner.CheckFile, partner.Upload)
	r.GET(Nft+"/:id", partner.GetNFTInfo)
	r.POST(Nickname, owner.SetNickname)
	r.GET(Nickname+"/:address", owner.GetNickname)

	var err error
	if gin.Mode() == gin.ReleaseMode {
		tool.Logger().Info("start gateway server", "port", port)
		err = r.RunTLS(":"+port, fullchain, privatekey)
	} else {
		err = r.Run(":" + port)
	}
	if err != nil {
		panic(err)
	}
}
