package main

import (
	"backend/api/gateway/owner"
	"github.com/gin-gonic/gin"
)

const (
	Auth = "/auth"
)

func main() {
	r := gin.Default()
	r.GET(Auth, owner.IsAuth)

	err := r.Run(":10321")
	if err != nil {
		panic(err)
	}
}
