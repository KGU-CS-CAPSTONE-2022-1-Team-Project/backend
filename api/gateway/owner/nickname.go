package owner

import (
	"backend/api/gateway/client"
	pb "backend/proto/owner"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type holderNickName struct {
	Address  string
	Nickname string
}

func SetNickname(ctx *gin.Context) {
	holder := &holderNickName{}
	err := ctx.Bind(holder)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "fail parsing json",
		})
		return
	}
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	req := pb.NicknameRequest{
		Address:  holder.Address,
		Nickname: holder.Nickname,
	}
	res, err := client.Owner().SetAnnoymousUser(timeout, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	switch res.Status.Code {
	case http.StatusBadRequest:
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "wrong params",
		})
	case http.StatusNotFound:
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "not found",
		})
	case http.StatusOK:
		ctx.JSON(http.StatusOK, gin.H{
			"message": "success",
		})

	case http.StatusForbidden:
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": res.Status.Message,
		})
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
	}
	return
}

func GetNickname(ctx *gin.Context) {
	address := ctx.Param("address")
	if address == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "empty",
		})
		return
	}
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	req := pb.NicknameRequest{
		Address: address,
	}
	res, err := client.Owner().GetAnnoymousUser(timeout, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	switch res.Status.Code {
	case http.StatusBadRequest:
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "wrong params",
		})
		return
	case http.StatusOK:
		ctx.JSON(http.StatusOK, gin.H{
			"message":  "success",
			"nickname": res.Nickname,
		})
		return
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
}
