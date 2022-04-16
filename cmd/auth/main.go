package main

import (
	"backend/api/auth"
	"github.com/gin-gonic/gin"
)

const (
	GOOGLE_AUTH          = "/auth/google"
	GOOGLE_AUTH_CALLBACK = "/auth/google/callback"
)

func main() {
	r := gin.Default()
	r.POST(GOOGLE_AUTH, auth.RequestResourceOwner)
	r.GET(GOOGLE_AUTH_CALLBACK, auth.GetTokenByGoogleServer)
	r.Run(":8000")
}
