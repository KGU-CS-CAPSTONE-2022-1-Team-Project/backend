package google

import (
	"backend/infrastructure/owner/dao"
	"backend/internal/owner"
	"backend/tool"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	userProfile "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"time"
)

var Config *oauth2.Config

func init() {
	if !viper.IsSet("web") {
		tool.ReadConfig("./configs/owner", "client_secret", "json")
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
	accessToken := owner.AccessToken{}
	err := owner.Validate(&accessToken, token)
	if err == nil {
		ctx.AbortWithStatusJSON(http.StatusOK, ResponseCommon{
			Message: "success",
		})
		return
	}
}

func GetUser(ctx *gin.Context) {
	headerAuth := ctx.Request.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(headerAuth, "Bearer ")
	tokenInfo := owner.AccessToken{}
	err := owner.GetAuthInfo(&tokenInfo, tokenString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ResponseCommon{Message: "파싱 실패"})
		return
	}
	tmp := dao.User{ID: tokenInfo.UserID}
	user, err := tmp.Read()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ResponseCommon{Message: "db조회 실패"})
		return
	}
	ctx.Set("user", user)
}

func RequestAuth(ctx *gin.Context) {
	url := Config.AuthCodeURL(
		ctx.ClientIP(),
		oauth2.AccessTypeOffline,
	)
	ctx.JSON(http.StatusTemporaryRedirect,
		ResponseAuth{
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
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ResponseCommon{
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
		ctx.AbortWithStatusJSON(http.StatusNotFound, ResponseCommon{Message: "이메일 권한 필요"})
	}
	user := owner.User{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Email:        email,
	}
	hasEmail := owner.User2UserDB(owner.User{Email: email})
	result, err := hasEmail.Read()
	var userDB dao.User
	if err == nil {
		result.AccessToken = token.AccessToken
		result.RefreshToken = token.RefreshToken
		err = result.Save()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError,
				ResponseCommon{Message: "갱신 실패"})
			return
		}
		userDB = *result
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		userDB = owner.User2UserDB(user)
		err = userDB.Create()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError,
				ResponseCommon{
					Message: "DB에러",
				})
			return
		}
	} else {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ResponseCommon{
			Message: "관리자에게 문의해주세요",
		})
		return
	}
	ctx.Set("userDB", userDB)
}

func CreateToken(ctx *gin.Context) {
	value, exist := ctx.Get("userDB")
	if !exist {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ResponseCommon{
			Message: "관리자에게 문의해주세요",
		})
		return
	}
	userDB := value.(dao.User)
	accessToken := owner.AccessToken{
		UserID:         userDB.ID,
		StandardClaims: jwt.StandardClaims{},
	}
	accessTokenString, err := owner.CreateTokenString(&accessToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseCommon{
				Message: "관리자에게 문의하세요",
			})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, ResponseAuth{
		Message:     "success",
		AccessToken: accessTokenString,
	})
}

func UpdateAddress(ctx *gin.Context) {
	tmp, exist := ctx.Get("user")
	if !exist {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ResponseCommon{Message: "유저정보교환 실패"})
		return
	}
	user := tmp.(*dao.User)
	if !user.IsAuthedStreamer {
		ctx.AbortWithStatusJSON(http.StatusForbidden, ResponseCommon{
			Message: "미인증된 채널",
		})
		return
	}
	holder := RequestAddr{}
	err := ctx.BindJSON(&holder)
	if err != nil || len(holder.Address) != 42 {
		log.Println(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseCommon{Message: "잘못된 파라미터"})
		return
	}
	user.Address = holder.Address
	err = user.Save()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ResponseCommon{Message: "업데이트 실패"})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, ResponseCommon{Message: "성공"})
	go RegisterContract(user.Address, user.ID)
}
