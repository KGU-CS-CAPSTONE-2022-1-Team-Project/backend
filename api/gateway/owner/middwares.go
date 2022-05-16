package owner

import (
	"backend/api/gateway/hosts"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type Result struct {
	Message     string `json:"message"`
	AuthUrl     string `json:"auth_url,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
}

var client = resty.New()

type holderAddress struct {
	Address string `json:"address"`
}

func IsAuth(ctx *gin.Context) {
	result := &Result{}
	client := resty.New()
	httpInfo, err := client.R().SetHeader("Authorization",
		ctx.Request.Header.Get("Authorization")).
		Post(hosts.Owner + "/owner/google")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "인증 서버 호출 실패",
		})
		return
	}
	err = json.Unmarshal(httpInfo.Body(), &result)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, Result{
			Message: "게이트웨이 서버 에러",
		})
		return
	}
	switch httpInfo.StatusCode() {
	case http.StatusTemporaryRedirect:
		ctx.Redirect(http.StatusFound, result.AuthUrl)
	case http.StatusOK:
		ctx.JSON(http.StatusOK, gin.H{
			"message": "validate",
		})
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, result)
	}
}

func Address(ctx *gin.Context) {
	holder := &holderAddress{}
	err := ctx.BindJSON(holder)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, Result{Message: "잘못된 파라미터"})
		return
	}
	httpInfo, err := client.R().SetHeader("Authorization",
		ctx.Request.Header.Get("Authorization")).
		SetBody(holder).
		Put(hosts.Owner + "/owner/address")
	result := &Result{}
	err = json.Unmarshal(httpInfo.Body(), &result)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, Result{
			Message: "인증 서버 호출 실패",
		})
		return
	}
	switch httpInfo.StatusCode() {
	case http.StatusOK:
		ctx.JSON(http.StatusOK, Result{
			Message: "success",
		})
	case http.StatusForbidden:
		ctx.JSON(http.StatusForbidden, result)
	case http.StatusBadRequest:
		ctx.JSON(http.StatusBadRequest, Result{Message: "서버에 값전달 실패"})
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, result)
		return
	}
}
