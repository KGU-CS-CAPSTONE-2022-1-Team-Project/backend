package auth

import (
	"backend/infrastructure/auth/dao"
	"crypto/rand"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"time"
)

var tokenSecret string

func init() {
	if !viper.IsSet("token_secert") {
		viper.SetConfigName("client_secret")
		viper.SetConfigType("json")
		viper.AddConfigPath("configs/auth")
		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("viper error: %v", err))
		}
	}
	tokenSecret = viper.GetString("token_secert")
}

type Token interface {
	validate(tokenString string) error
	create() (string, error)
	parser(tokenString string) error
}

func Validate(t Token, tokenString string) error {
	return t.validate(tokenString)
}

func CreateTokenString(t Token) (string, error) {
	return t.create()
}

func GetAuthInfo(t Token, tokenString string) error {
	return t.parser(tokenString)
}

const otpChars = "0123456789"

// MakeTokenId
// Reference: https://stackoverflow.com/questions/39481826/generate-6-digit-verification-code-with-golang
func MakeTokenId() (string, error) {
	buffer := make([]byte, 9)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", errors.Wrap(err, "tokenId생성 시")
	}
	otpCharsLength := len(otpChars)
	for i := 0; i < 9; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}
	return string(buffer), nil
}

type AccessToken struct {
	UserID string `json:"UserID"`
	jwt.StandardClaims
}

type RefreshToken struct {
	AccessTokenString string `json:"accessToken"`
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
		user := dao.User{
			ID:              r.UserID,
			TokenIdentifier: r.Id,
		}
		if _, err := user.Read(); err != nil {
			return nil, errors.Wrap(err, "db 조회 실패")
		}
		return []byte(tokenSecret), nil
	})
	return err
}

func (r *RefreshToken) validate(tokenString string) error {
	_, err := jwt.ParseWithClaims(tokenString, r, func(token *jwt.Token) (interface{}, error) {
		// 1. 해싱알고리즘
		if tokenMethod, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || tokenMethod != jwt.SigningMethodHS256 {
			return nil, errors.Wrap(ErrValidateToken, "잘못된 알고리즘")
		}
		// 2. standard기준 확인(exp,iat)
		if err := r.Valid(); err != nil {
			return nil, errors.Wrap(err, "토큰 유효기간 혹은 생성시간")
		}
		// 3. 액세스 토큰 내의 uid확인
		accessToken := AccessToken{}
		_, err := jwt.ParseWithClaims(r.AccessTokenString, &accessToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		})
		if err != nil {
			return nil, ErrDecrypt
		}
		user := dao.User{
			ID: accessToken.UserID,
		}

		if _, err = user.Read(); err != nil {
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

func (r *RefreshToken) create() (string, error) {
	duration := 14 * 24 * time.Hour
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

func (r *RefreshToken) parser(tokenString string) error {
	_, err := jwt.ParseWithClaims(tokenString, r, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	return err
}
