package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type ResultAuth struct {
	State        int    `json:"state"`
	AuthUrl      string `json:"auth_url,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func IsAuth(ctx *gin.Context) {
	result := &ResultAuth{}
	client := resty.New()
	_, err := client.R().SetResult(result).Post("http://localhost:8000/auth/google")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "인증 서버 에러",
		})
		return
	}
	switch result.State {
	case 1:
		ctx.Redirect(http.StatusMovedPermanently, result.AuthUrl)
	case 2:
		ctx.JSON(http.StatusOK, gin.H{
			"message":       "success",
			"access_token":  result.AccessToken,
			"refresh_token": result.RefreshToken,
		})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "인증서버 에러",
		})
	}
}
