package test

import (
	"backend/infrastructure/owner/dao"
	"backend/internal/owner"
	"backend/internal/owner/youtuber"
	"backend/tool"
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
	tool.ReadConfig("./configs/owner", "test_secret", "json")
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
	mock := owner.Owner{
		Email: "test@test.com",
	}
	userDB := dao.Owner{
		Email: mock.Email,
	}
	defer func() {
		dbConnection().Unscoped().
			Delete(&userDB)
	}()
	err := userDB.Create()
	require.Nil(t, err, "생성 실패")
	result, err := userDB.Read()
	require.Nil(t, err, "조회 실패")
	assert.Equal(t, mock.Email, result.Email, "잘못된 입력")
}

func TestMigration(t *testing.T) {
	temp1 := dao.Owner{}
	temp2 := dao.User{}
	err := temp1.Migration()
	require.Nil(t, err, "Owner 마이그레이션", err)
	err = temp2.Migration()
	require.Nil(t, err, "Owner 마이그레이션", err)
}

func TestAccessToken(t *testing.T) {
	userDB := dao.Owner{
		Email: "test@test.com",
	}
	defer func() {
		dbConnection().Unscoped().
			Delete(&userDB)
	}()
	err := userDB.Create()
	require.Nil(t, err, "생성 실패")
	result, err := userDB.Read()
	src := owner.AccessToken{
		UserID:         result.ID,
		StandardClaims: jwt.StandardClaims{},
	}
	accessTokenString, err := owner.CreateTokenString(&src)
	assert.Nil(t, err, "생성 실패", err)
	dst := owner.AccessToken{}
	err = owner.TokenValidate(&dst, accessTokenString)
	assert.Nil(t, err, "정상토큰을 오류로 인식", err, errors.Cause(err))
	assert.True(t, reflect.DeepEqual(src, dst), "같은 토큰이나 서로 다른정보로 인식")

	// Fail cases
	notFoundUUIDToken := owner.AccessToken{
		UserID:         uuid.NewString(),
		StandardClaims: jwt.StandardClaims{},
	}
	err = owner.TokenValidate(&dst, accessTokenString)
	assert.Nil(t, err, "정상토큰을 검증실패", "에러내용: ", err)
	accessTokenString, _ = owner.CreateTokenString(&notFoundUUIDToken)
	err = owner.TokenValidate(&notFoundUUIDToken, accessTokenString)
	assert.NotNil(t, err, "db에 존재하지 않는 토큰")
	expiredToken := owner.AccessToken{
		UserID: result.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(-10 * time.Second).Unix(),
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, &expiredToken)
	accessTokenString, _ = at.SignedString([]byte(tokenSecret))
	err = owner.TokenValidate(&notFoundUUIDToken, accessTokenString)
	assert.NotNil(t, err, "만료된 토큰")
	wrongAlgorithmToken := owner.AccessToken{
		UserID: result.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(1 * time.Hour).Unix(),
			IssuedAt:  time.Now().UTC().Unix(),
		},
	}
	at = jwt.NewWithClaims(jwt.SigningMethodHS384, &wrongAlgorithmToken)
	accessTokenString, err = at.SignedString([]byte(tokenSecret))
	require.Nil(t, err, "다른알고리즘생성실패")
	err = owner.TokenValidate(&wrongAlgorithmToken, accessTokenString)
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

func TestValidateAddrAndNickname(t *testing.T) {
	conn := dbConnection()
	usr := owner.User{
		Address:  "0xd753883f95e059abbde245c150cb87dc513044f7",
		Nickname: "wahaha",
	}
	defer conn.Delete(usr)
	err := owner.Validate(&usr)
	require.Nil(t, err, "정상값을 에러처리", err)
	userDAO := owner.User2UserDB(usr)
	err = userDAO.Save()
	assert.Nil(t, err, "정상값 저장 실패", err)

	failCases := [4]owner.User{}
	for caseIdx, failCase := range failCases {
		switch caseIdx {
		case 0:
			// 빈 값
			failCase = owner.User{}
			err = owner.Validate(&failCase)
		case 1:
			// 잘못된 주소값
			failCase = owner.User{
				Address:  "0xd753883f95e059abbde245c150cb87dc51304",
				Nickname: "wahaha",
			}
		case 2:
			// 짧은 닉네임길이
			failCase = owner.User{
				Address:  "0xd753883f95e059abbde245c150cb87dc51304",
				Nickname: "김김김",
			}
		case 3:
			// 긴 닉네임길이
			failCase = owner.User{
				Address:  "0xd753883f95e059abbde245c150cb87dc51304",
				Nickname: "김김김김김김김김김김김김",
			}
		}
		err = owner.Validate(&failCase)
		assert.NotNil(t, err, "에러값을 정상처리", "case:", caseIdx)
	}
}
