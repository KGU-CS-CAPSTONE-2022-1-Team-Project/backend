package test

import (
	"backend/internal/auth"
	"backend/internal/auth/dao"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"reflect"
	"testing"
	"time"
)

var tokenSecret string

func init() {
	viper.SetConfigName("client_secret")
	viper.SetConfigType("json")
	viper.SetConfigType("")
	viper.AddConfigPath("configs/auth")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("viper error: %v", err))
	}
	tokenSecret = viper.GetString("token_secert")
}

func dbconnection() *gorm.DB {
	config := viper.GetStringMapString("db")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config["user"], config["password"], config["host"], config["port"], config["main_db"])
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func TestCreateUser(t *testing.T) {
	user := auth.User{
		Email: "test@test.com",
	}
	userDB := dao.User{
		Email: user.Email,
	}
	defer func() {
		dbconnection().Unscoped().
			Delete(&userDB)
	}()
	err := userDB.Create()
	require.Nil(t, err, "생성 실패")
	result, err := userDB.Read()
	require.Nil(t, err, "조회 실패")
	assert.Equal(t, user.Email, result.Email, "잘못된 입력")
}

func TestMigration(t *testing.T) {
	user := dao.User{}
	err := user.Migration()
	require.Nil(t, err, "마이그레이션")
}

func TestAccessToken(t *testing.T) {
	tokenID, err := auth.MakeTokenId()
	require.Nil(t, err, "토큰id생성실패")
	userDB := dao.User{
		Email:           "test@test.com",
		TokenIdentifier: tokenID,
	}
	defer func() {
		dbconnection().Unscoped().
			Delete(&userDB)
	}()
	err = userDB.Create()
	require.Nil(t, err, "생성 실패")
	user, err := userDB.Read()
	src := auth.AccessToken{
		UserID:         user.ID,
		StandardClaims: jwt.StandardClaims{Id: tokenID},
	}
	accessTokenString, err := auth.CreateTokenString(&src)
	assert.Nil(t, err, "생성 실패", err)
	dst := auth.AccessToken{}
	err = auth.Validate(&dst, accessTokenString)
	assert.Nil(t, err, "정상토큰을 오류로 인식", err, errors.Cause(err))
	assert.True(t, reflect.DeepEqual(src, dst), "같은 토큰이나 서로 다른정보로 인식")

	// Fail cases
	notFoundUUIDToken := auth.AccessToken{
		UserID:         uuid.NewString(),
		StandardClaims: jwt.StandardClaims{},
	}
	err = auth.Validate(&dst, accessTokenString)
	assert.Nil(t, err, "정상토큰을 검증실패", "에러내용: ", err)
	accessTokenString, _ = auth.CreateTokenString(&notFoundUUIDToken)
	err = auth.Validate(&notFoundUUIDToken, accessTokenString)
	assert.NotNil(t, err, "db에 존재하지 않는 토큰")
	expiredToken := auth.AccessToken{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			Id:        tokenID,
			ExpiresAt: time.Now().Add(-10 * time.Second).Unix(),
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, &expiredToken)
	accessTokenString, _ = at.SignedString([]byte(tokenSecret))
	err = auth.Validate(&notFoundUUIDToken, accessTokenString)
	assert.NotNil(t, err, "만료된 토큰")
	wrongAlgorithmToken := auth.AccessToken{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			Id:        tokenID,
			ExpiresAt: time.Now().UTC().Add(1 * time.Hour).Unix(),
			IssuedAt:  time.Now().UTC().Unix(),
		},
	}
	at = jwt.NewWithClaims(jwt.SigningMethodHS384, &wrongAlgorithmToken)
	accessTokenString, err = at.SignedString([]byte(tokenSecret))
	require.Nil(t, err, "다른알고리즘생성실패")
	err = auth.Validate(&wrongAlgorithmToken, accessTokenString)
	assert.NotNil(t, err, "잘못된 알고리즘 통과")
}

func TestRefreshToken(t *testing.T) {
	tokenID, err := auth.MakeTokenId()
	require.Nil(t, err, "토큰id생성실패")
	userDB := dao.User{
		Email:           "test@test.com",
		TokenIdentifier: tokenID,
	}
	defer func() {
		dbconnection().Unscoped().
			Delete(&userDB)
	}()
	err = userDB.Create()
	require.Nil(t, err, "생성 실패")
	user, err := userDB.Read()
	accessToken := auth.AccessToken{
		UserID:         user.ID,
		StandardClaims: jwt.StandardClaims{Id: tokenID},
	}
	accessTokenString, err := auth.CreateTokenString(&accessToken)
	refreshToken := auth.RefreshToken{
		AccessTokenString: accessTokenString,
		StandardClaims: jwt.StandardClaims{
			Id: tokenID,
		},
	}
	refreshTokenString, err := auth.CreateTokenString(&refreshToken)
	refreshToken2 := auth.RefreshToken{}
	err = auth.GetAuthInfo(&refreshToken2, refreshTokenString)
	assert.Nil(t, err, "파싱에러발생", err)
	assert.True(t, reflect.DeepEqual(refreshToken, refreshToken2), "같은토큰들, 다른값")
	err = auth.Validate(&refreshToken, refreshTokenString)
	assert.Nil(t, err, "정상토큰 검증에러 발생", err)

	// fail case
	wrongAlgorithm := auth.RefreshToken{
		AccessTokenString: accessTokenString,
		StandardClaims: jwt.StandardClaims{
			Id:        tokenID,
			ExpiresAt: time.Now().Add(10 * time.Second).Unix(),
		},
	}
	refreshTokenString, err = jwt.NewWithClaims(jwt.SigningMethodHS384, &wrongAlgorithm).
		SignedString([]byte(tokenSecret))
	err = auth.Validate(&refreshToken, refreshTokenString)
	assert.NotNil(t, err, "잘못된 알고리즘")

	expireToken := auth.RefreshToken{
		AccessTokenString: accessTokenString,
		StandardClaims: jwt.StandardClaims{
			Id:        tokenID,
			ExpiresAt: time.Now().Add(-10 * time.Second).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	refreshTokenString, err = jwt.NewWithClaims(jwt.SigningMethodHS256, &expireToken).
		SignedString([]byte(tokenSecret))
	err = auth.Validate(&refreshToken, tokenSecret)
	assert.NotNil(t, err, "만료된 토큰")
	accessToken = auth.AccessToken{UserID: uuid.NewString()}
	accessTokenString, err = auth.CreateTokenString(&accessToken)
	require.Nil(t, err, "가짜 accessToken생성 실패")
	hasNotDBToken := auth.RefreshToken{
		AccessTokenString: accessTokenString,
		StandardClaims:    jwt.StandardClaims{Id: "123"},
	}
	refreshTokenString, err = auth.CreateTokenString(&hasNotDBToken)
	require.Nil(t, err, "토큰 생성 실패", err, errors.Cause(err))
	err = auth.Validate(&hasNotDBToken, refreshTokenString)
	assert.NotNil(t, err, "db에 존재하지 않는 토큰")
}
