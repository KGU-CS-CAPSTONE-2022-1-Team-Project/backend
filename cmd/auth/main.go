package main

import (
	"backend/api/auth"
	"github.com/gin-gonic/gin"
)

const (
	GOOGLE_AUTH          = "/auth/google"
	GOOGLE_AUTH_CALLBACK = "/auth/google/callback"
	REFRESH_TOKEN        = "/auth/refresh"
)

func main() {
	r := gin.Default()
	r.POST(GOOGLE_AUTH, auth.CheckUser, auth.RequestAuth)
	r.GET(GOOGLE_AUTH_CALLBACK, auth.GetTokenByGoogleServer, auth.RegisterUser, auth.CreateToken)
	r.PATCH(REFRESH_TOKEN, auth.CheckRefresh, auth.CreateToken)
	err := r.Run(":8000")
	if err != nil {
		panic(err)
	}
}
