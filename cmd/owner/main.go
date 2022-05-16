package main

import (
	"backend/api/owner/google"
	"github.com/gin-gonic/gin"
)

const (
	GoogleAuth         = "/owner/google"
	GoogleAuthCallback = "/owner/google/callback"
	Youtuber           = "/owner/youtuber"
	Address            = "/owner/address"
)

func main() {
	r := gin.Default()
	r.POST(GoogleAuth, google.CheckNotUser, google.RequestAuth)
	r.GET(GoogleAuthCallback,
		google.GetTokenByGoogleServer,
		google.RegisterUser,
		google.CreateToken)
	r.PUT(Youtuber, google.GetUser, google.ISYoutuber, google.SetChannel)
	r.PUT(Address, google.GetUser, google.UpdateAddress)
	err := r.Run(":8000")
	if err != nil {
		panic(err)
	}
}
