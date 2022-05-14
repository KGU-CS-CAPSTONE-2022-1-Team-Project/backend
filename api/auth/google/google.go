package google

import (
	"backend/infrastructure/auth/dao"
	"backend/internal/auth"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	userProfile "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

type Response struct {
	Message      string `json:"message"`
	AuthUrl      string `json:"auth_url,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type RequestRefresh struct {
	RefreshToken string `json:"refresh_token"`
}

type holderAddress struct {
	Address string `json:"address"`
}

var Config *oauth2.Config

func init() {
	if !viper.IsSet("web") {
		viper.SetConfigName("client_secret")
		viper.SetConfigType("json")
		viper.AddConfigPath("./configs/auth")
		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("viper error: %v", err))
		}
	}
	infoAuth := viper.GetStringMapStringSlice("web")
	Config = &oauth2.Config{
		ClientID:     infoAuth["client_id"][0],
		ClientSecret: infoAuth["client_secret"][0],
		RedirectURL:  infoAuth["redirect_uris"][0],
		Scopes:       infoAuth["scopes"],
		Endpoint:     google.Endpoint,
	}
}

func getGoogleEmail(ctx context.Context, token *oauth2.Token) (string, error) {
	source := Config.TokenSource(ctx, token)
	client, err := userProfile.NewService(ctx, option.WithTokenSource(source))
	userInfo, err := client.Userinfo.Get().Do()
	if err != nil {
		return "", err
	}
	return userInfo.Email, nil
}

func CheckNotUser(ctx *gin.Context) {
	headerAuth := ctx.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(headerAuth, "Bearer ")
	if token == "" {
		ctx.Next()
		return
	}
	accessToken := auth.AccessToken{}
	err := auth.Validate(&accessToken, token)
	if err == nil {
		ctx.AbortWithStatusJSON(http.StatusOK, Response{
			Message: "success",
		})
		return
	}
}

func GetUser(ctx *gin.Context) {
	headerAuth := ctx.Request.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(headerAuth, "Bearer ")
	tokenInfo := auth.AccessToken{}
	err := auth.GetAuthInfo(&tokenInfo, tokenString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: "파싱 실패"})
		return
	}
	tmp := dao.User{ID: tokenInfo.UserID, TokenIdentifier: tokenInfo.Id}
	user, err := tmp.Read()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: "db조회 실패"})
		return
	}
	ctx.Set("user", user)
}

func CheckRefresh(ctx *gin.Context) {
	requsetInfo := RequestRefresh{}
	err := ctx.BindJSON(&requsetInfo)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, Response{Message: "파라미터 조회 실패"})
		return
	}
	refreshToken := auth.RefreshToken{}
	err = auth.Validate(&refreshToken, requsetInfo.RefreshToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, Response{
			Message: "유효하지 않는 토큰",
		})
	}
	accessToken := auth.AccessToken{}
	err = auth.GetAuthInfo(&accessToken, refreshToken.AccessTokenString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{
			Message: "만료되었으면서",
		})
	}
	searchValues := dao.User{ID: accessToken.UserID}
	userDB, err := searchValues.Read()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{
			Message: "관리자에게 문의하세요"})
	}
	ctx.Set("userDB", *userDB)
}

func RequestAuth(ctx *gin.Context) {
	url := Config.AuthCodeURL(
		ctx.ClientIP(),
		oauth2.AccessTypeOffline,
	)
	ctx.JSON(http.StatusTemporaryRedirect,
		Response{
			Message: "인증 필요",
			AuthUrl: url,
		})
}

func GetTokenByGoogleServer(ctx *gin.Context) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	code := ctx.Query("code")
	token, err := Config.Exchange(ctxTimeout, code)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{
			Message: "exchange실패",
		})
		return
	}
	ctx.Set("token", token)
}

func RegisterUser(ctx *gin.Context) {
	tokenType, _ := ctx.Get("token")
	token := tokenType.(*oauth2.Token)
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	email, err := getGoogleEmail(ctxTimeout, token)
	if err != nil || email == "" {
		ctx.AbortWithStatusJSON(http.StatusNotFound, Response{Message: "이메일 권한 필요"})
	}

	tokenID, err := auth.MakeTokenId()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: "토큰id생성실패"})
	}
	user := auth.User{
		TokenIdentifier: tokenID,
		AccessToken:     token.AccessToken,
		RefreshToken:    token.RefreshToken,
		Email:           email,
	}
	hasEmail := auth.User2UserDB(auth.User{Email: email})
	result, err := hasEmail.Read()
	var userDB dao.User
	if err == nil {
		result.AccessToken = token.AccessToken
		result.RefreshToken = token.RefreshToken
		result.TokenIdentifier = tokenID
		err = result.Save()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError,
				Response{Message: "갱신 실패"})
			return
		}
		userDB = *result
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		userDB = auth.User2UserDB(user)
		err = userDB.Create()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError,
				Response{
					Message: "DB에러",
				})
			return
		}
	} else {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{
			Message: "관리자에게 문의해주세요",
		})
		return
	}
	ctx.Set("userDB", userDB)
}

func CreateToken(ctx *gin.Context) {
	value, exist := ctx.Get("userDB")
	if !exist {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{
			Message: "관리자에게 문의해주세요",
		})
		return
	}
	userDB := value.(dao.User)
	accessToken := auth.AccessToken{
		UserID:         userDB.ID,
		StandardClaims: jwt.StandardClaims{Id: userDB.TokenIdentifier},
	}
	accessTokenString, err := auth.CreateTokenString(&accessToken)
	refreshToken := auth.RefreshToken{
		AccessTokenString: accessTokenString,
		StandardClaims: jwt.StandardClaims{
			Id: accessToken.Id,
		},
	}
	refreshTokenString, err := auth.CreateTokenString(&refreshToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			Response{
				Message: "관리자에게 문의하세요",
			})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, Response{
		Message:      "success",
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	})
}

func UpdateAddress(ctx *gin.Context) {
	tmp, exist := ctx.Get("user")
	if !exist {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: "유저정보교환 실패"})
		return
	}
	user := tmp.(*dao.User)
	if !user.IsAuthedStreamer {
		ctx.AbortWithStatusJSON(http.StatusForbidden, Response{
			Message: "미인증된 채널",
		})
		return
	}
	holder := &holderAddress{}
	err := ctx.BindJSON(holder)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, Response{Message: "잘못된 파라미터"})
		return
	}
	user.Address = holder.Address
	err = user.Save()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: "업데이트 실패"})
		return
	}
	ctx.JSON(http.StatusOK, Response{Message: "성공"})
}
