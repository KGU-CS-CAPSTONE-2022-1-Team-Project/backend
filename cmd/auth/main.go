package main

import (
	"backend/api/auth/google"
	"github.com/gin-gonic/gin"
)

const (
	GoogleAuth         = "/auth/google"
	GoogleAuthCallback = "/auth/google/callback"
	Youtuber           = "/auth/authorship/youtuber"
	Address            = "/auth/authorship/youtuber/address"
)

func main() {
	r := gin.Default()
	r.POST(GoogleAuth, google.CheckNotUser, google.RequestAuth)
	r.GET(GoogleAuthCallback, google.GetTokenByGoogleServer, google.RegisterUser, google.CreateToken)

	r.PUT(Youtuber, google.GetUser, google.ISYoutuber)
	r.PUT(Address, google.GetUser, google.UpdateAddress)

	err := r.Run(":8000")
	if err != nil {
		panic(err)
	}
}
