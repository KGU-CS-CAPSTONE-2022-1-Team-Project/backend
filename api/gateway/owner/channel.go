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
	switch res.Status.Code {
	case http.StatusNotFound:
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "not exists channel",
		})
	case http.StatusOK:
		ctx.JSON(http.StatusOK, gin.H{
			"title":                   res.Title,
			"description":             res.Description,
			"image":                   res.Image,
			"external_link":           res.Url,
			"seller_fee_basis_points": 1000,
			"fee_recipient":           res.Address,
		})
	}
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
}
