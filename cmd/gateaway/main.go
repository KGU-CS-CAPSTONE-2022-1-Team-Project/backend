package main

import (
	"backend/api/gateway"
	"github.com/gin-gonic/gin"
)

const (
	Auth = "/auth"
)

func main() {
	r := gin.Default()
	r.GET(Auth, gateway.IsAuth)
	r.Run(":10321")
}
