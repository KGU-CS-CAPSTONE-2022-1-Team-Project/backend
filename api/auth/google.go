package auth

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"time"
)

func Auth() {
	config := &oauth2.Config{
		ClientID:     viper.GetString("ClienntID"),
		ClientSecret: viper.GetString("ClientSecret"),
		RedirectURL:  viper.GetString("RedirectUri"),
	}
	fmt.Println(config)
}

const (
	REQUEST_AUTH = 1
	TOKEN        = iota
)

var config *oauth2.Config

func init() {
	viper.SetConfigName("client_secret")
	viper.SetConfigType("json")
	viper.AddConfigPath("./configs/auth")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("viper error: %v", err))
	}
	infoAuth := viper.GetStringMapStringSlice("web")
	config = &oauth2.Config{
		ClientID:     infoAuth["client_id"][0],
		ClientSecret: infoAuth["client_secret"][0],
		RedirectURL:  infoAuth["redirect_uris"][0],
		Scopes:       infoAuth["scopes"],
		Endpoint:     google.Endpoint,
	}
}

func RequestResourceOwner(ctx *gin.Context) {
	url := config.AuthCodeURL(
		ctx.ClientIP(),
		oauth2.AccessTypeOffline,
	)
	ctx.JSON(http.StatusOK, gin.H{
		"state":    REQUEST_AUTH,
		"auth_url": url,
	})
}
func GetTokenByGoogleServer(ctx *gin.Context) {
	ctxTimeout, _ := context.WithTimeout(context.Background(), 5*time.Second)
	code := ctx.Query("code")
	token, err := config.Exchange(ctxTimeout, code)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "exchange실패",
		})
		return
	}
	fmt.Println("access token:", token.AccessToken)
	fmt.Println("refresh token:", token.RefreshToken)
	fmt.Println("token time", token.Expiry)
	fmt.Println("token type", token.TokenType)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
