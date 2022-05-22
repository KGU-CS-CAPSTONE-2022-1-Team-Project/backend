package owner

import (
	"backend/api/gateway/client"
	pb "backend/proto/owner"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type holderAddress struct {
	Address string
}

func GetChannel(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "ID not found",
		})
		return
	}
	timeout, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	req := pb.ChannelRequest{Id: id}
	res, err := client.Owner().GetChannel(timeout, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "내부서버 오류",
		})
		return
	}
	if res.IsEmpty {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "존재하지 않는 채널",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"title":       res.Title,
		"description": res.Description,
		"image":       res.Image,
		"url":         res.Url,
	})
}
