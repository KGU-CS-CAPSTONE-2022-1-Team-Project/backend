package test

import (
	"backend/infrastructure/auth/dao"
	"backend/internal/auth"
	"backend/internal/auth/youtuber"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"reflect"
	"testing"
	"time"
)

var tokenSecret string

func init() {
	viper.SetConfigName("test_secret")
	viper.SetConfigType("json")
	viper.AddConfigPath("../configs/auth")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("viper error: %v", err))
	}
	tokenSecret = viper.GetString("token_secert")
}

func dbConnection() *gorm.DB {
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
		dbConnection().Unscoped().
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
		dbConnection().Unscoped().
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

func TestYoutubeApi(t *testing.T) {
	require.True(t, viper.IsSet("web"), "설정 조회 실패")
	require.True(t, viper.IsSet("access_token"), "테스트용 설정파일 조회 실패")
	require.True(t, viper.IsSet("refresh_token"), "테스트용 설정파일 조회 실패")
	infoAuth := viper.GetStringMapStringSlice("web")
	config := &oauth2.Config{
		ClientID:     infoAuth["client_id"][0],
		ClientSecret: infoAuth["client_secret"][0],
		RedirectURL:  infoAuth["redirect_uris"][0],
		Scopes:       infoAuth["scopes"],
		Endpoint:     google.Endpoint,
	}
	token := oauth2.Token{
		AccessToken:  viper.GetString("access_token"),
		TokenType:    "Bearer",
		RefreshToken: viper.GetString("refresh_token"),
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	source := config.TokenSource(ctx, &token)
	youtubeService, err := youtube.NewService(ctx, option.WithTokenSource(source))
	require.Nil(t, err, "유튜브 서비스에러", err)
	result, err := youtubeService.Channels.List(
		[]string{"auditDetails", "statistics", "snippet"},
	).Mine(true).Do()
	require.Nil(t, err, "유튜브 채널 api에러", err)
	assert.Equal(t, len(result.Items), 1, "1개 초과 혹은0개의 결과값", len(result.Items))
}

func TestValidateYoutuber(t *testing.T) {
	validateMock := youtube.Channel{
		AuditDetails: &youtube.ChannelAuditDetails{
			CommunityGuidelinesGoodStanding: true,
			ContentIdClaimsGoodStanding:     true,
			CopyrightStrikesGoodStanding:    true,
		},
		ContentOwnerDetails: &youtube.ChannelContentOwnerDetails{
			ContentOwner: "owner",
			TimeLinked:   "time",
		},
		ConversionPings: nil,
		Statistics: &youtube.ChannelStatistics{
			CommentCount:          youtuber.MinCommentCount,
			HiddenSubscriberCount: false,
			SubscriberCount:       youtuber.MinSubscriber,
			ViewCount:             youtuber.MinViewerCount,
			VideoCount:            youtuber.MinVideoCount,
		},
	}
	validateMock.HTTPStatusCode = http.StatusOK
	err := youtuber.ValidateChannel(&validateMock)
	assert.Nil(t, err, "유효한 채널", err)
	invalidateMocks := [6]youtube.Channel{}
	for idx := range invalidateMocks {
		var buffer bytes.Buffer
		var dst *youtube.Channel
		encoder := json.NewEncoder(&buffer)
		decoder := json.NewDecoder(&buffer)
		err = encoder.Encode(validateMock)
		require.Nil(t, err, "deep copy준비문제")
		err = decoder.Decode(&dst)
		require.Nil(t, err, "deep copy준비문제")
		invalidateMocks[idx] = *dst
	}
	for idx, mock := range invalidateMocks {
		switch idx {
		case 0:
			mock.AuditDetails.CopyrightStrikesGoodStanding = false
			err = youtuber.ValidateChannel(&mock)
			assert.NotNil(t, err, "저작권문제가 있는 채널")
		case 1:
			mock.ContentOwnerDetails = nil
			err = youtuber.ValidateChannel(&mock)
			assert.NotNil(t, err, "파트너가 아닌 채널")
		case 2:
			mock.Statistics.ViewCount = 0
			err = youtuber.ValidateChannel(&mock)
			assert.NotNil(t, err, "조회수 부족")
		case 3:
			mock.Statistics.SubscriberCount = 0
			err = youtuber.ValidateChannel(&mock)
			assert.NotNil(t, err, "구독자수 부족")
		case 4:
			mock.Statistics.CommentCount = 0
			err = youtuber.ValidateChannel(&mock)
			assert.NotNil(t, err, "댓글수 부족")
		case 5:
			mock.Statistics.VideoCount = 0
			err = youtuber.ValidateChannel(&mock)
			assert.NotNil(t, err, "영상 개수 부족")
		}
	}
}
