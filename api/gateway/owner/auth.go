package owner

import (
	"backend/api/gateway/client"
	pb "backend/proto/owner"
	"backend/tool"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"strings"
	"time"
)

var redirectUri string

func init() {
	if redirectUri == "" {
		tool.ReadConfig("./configs/gateway", "client_secret", "json")
		redirectUri = viper.GetString("redirect_uri")
	}
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
	ctx.Redirect(http.StatusTemporaryRedirect, redirectUri+res.AccessToken)
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
