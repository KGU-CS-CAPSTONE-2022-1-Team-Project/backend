package owner

import (
	"backend/api/gateway/client"
	pb "backend/proto/owner"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type holderAddress struct {
	Address string
}

func Redirecting(ctx *gin.Context) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	res, err := client.Owner().Google(timeout, &pb.LoginRequest{Ip: ctx.ClientIP()})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "내부서버 오류",
		})
	}
	ctx.Redirect(http.StatusTemporaryRedirect, res.AuthUrl)
}

func RegisterUser(ctx *gin.Context) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	code := ctx.Query("code")
	req := pb.RegisterRequest{Code: code}
	res, err := client.Owner().GoogleCallBack(timeout, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "내부서버 오류",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"access_token": res.AccessToken,
	})
}

func AuthYoutuber(ctx *gin.Context) {
	auth := ctx.Request.Header.Get("Authorization")
	if auth == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Not Find Authorization Header",
		})
		return
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == auth {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Could not find token",
		})
		return
	}
	holder := holderAddress{}
	err := ctx.Bind(&holder)
	if err != nil || holder.Address == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Could not find address",
		})
	}

	req := pb.AddressRequest{Address: holder.Address, AccessToken: token}
	timeout, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	res, err := client.Owner().SaveAddress(timeout, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "내부 서버 에러",
		})
		return
	}
	if !res.IsValidate {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "인증 실패" + err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
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
