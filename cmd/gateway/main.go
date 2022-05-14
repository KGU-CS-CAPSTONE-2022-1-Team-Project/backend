package main

import (
	"backend/api/gateway/auth"
	"github.com/gin-gonic/gin"
)

const (
	Auth = "/auth"
)

func main() {
	r := gin.Default()
	r.GET(Auth, auth.IsAuth)
	err := r.Run(":10321")
	if err != nil {
		panic(err)
	}
}
