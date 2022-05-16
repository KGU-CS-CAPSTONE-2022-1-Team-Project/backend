package main

import (
	"backend/api/owner/google"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	GoogleAuth         = "/owner/google"
	GoogleAuthCallback = "/owner/google/callback"
	Youtuber           = "/owner/youtuber"
	Address            = "/owner/address"
)

var port string

func init() {
	if !viper.IsSet("port") {
		viper.SetConfigName("client_secert")
		viper.SetConfigType("json")
		viper.AddConfigPath("configs/owner")
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
	}
	port = viper.GetString("port")
}

func main() {
	r := gin.Default()
	r.POST(GoogleAuth, google.CheckNotUser, google.RequestAuth)
	r.GET(GoogleAuthCallback,
		google.GetTokenByGoogleServer,
		google.RegisterUser,
		google.CreateToken)
	r.PUT(Youtuber, google.GetUser, google.ISYoutuber, google.SetChannel)
	r.PUT(Address, google.GetUser, google.UpdateAddress)
	r.GET(Youtuber+"/:id", google.GetChannelInfo)
	err := r.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
