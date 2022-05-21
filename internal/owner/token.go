package owner

import (
	"backend/infrastructure/owner/dao"
	"backend/tool"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"time"
)

var tokenSecret string

func init() {
	if !viper.IsSet("token_secert") {
		tool.ReadConfig("./configs/owner", "client_secret", "json")
	}
	tokenSecret = viper.GetString("token_secert")
}

type Token interface {
	validate(tokenString string) error
	create() (string, error)
	parser(tokenString string) error
}

func TokenValidate(t Token, tokenString string) error {
	return t.validate(tokenString)
}

func CreateTokenString(t Token) (string, error) {
	return t.create()
}

func GetAuthInfo(t Token, tokenString string) error {
	return t.parser(tokenString)
}

type AccessToken struct {
	UserID string `json:"UserID"`
	jwt.StandardClaims
}

func (r *AccessToken) validate(tokenString string) error {
	_, err := jwt.ParseWithClaims(tokenString, r, func(token *jwt.Token) (interface{}, error) {
		// 1. 해싱알고리즘
		if tokenMethod, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || tokenMethod != jwt.SigningMethodHS256 {
			return nil, errors.Wrap(ErrValidateToken, "잘못된 알고리즘")
		}
		// 2. standard기준 확인(exp,iat)
		if err := r.Valid(); err != nil {
			return nil, errors.Wrap(err, "토큰 유효기간 혹은 생성시간")
		}
		// 3. db의 데이터와 일치하는지확인
		user := dao.Owner{
			ID: r.UserID,
		}
		if _, err := user.Read(); err != nil {
			return nil, errors.Wrap(err, "db 조회 실패")
		}
		return []byte(tokenSecret), nil
	})
	return err
}

func (r *AccessToken) create() (string, error) {
	duration := 1 * time.Hour
	r.ExpiresAt = time.Now().UTC().Add(duration).Unix()
	r.IssuedAt = time.Now().UTC().Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, r)
	return at.SignedString([]byte(tokenSecret))
}

func (r *AccessToken) parser(tokenString string) error {
	_, err := jwt.ParseWithClaims(tokenString, r, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	return err
}
