package auth

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type ResultAuth struct {
	Message      string `json:"message"`
	AuthUrl      string `json:"auth_url,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func IsAuth(ctx *gin.Context) {
	result := &ResultAuth{}
	client := resty.New()
	httpInfo, err := client.R().SetHeader("Authorization", ctx.Request.Header.Get("Authorization")).Post("http://localhost:8000/auth/google")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "인증 서버 에러",
		})
		return
	}
	err = json.Unmarshal(httpInfo.Body(), &result)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "게이트웨이 서버 에러",
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
		return
	}
}
