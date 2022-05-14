package google

import (
	"backend/infrastructure/auth/dao"
	"backend/internal/auth/youtuber"
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"net/http"
	"time"
)

func ISYoutuber(ctx *gin.Context) {
	tmp, exist := ctx.Get("user")
	if !exist {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: "유저정보교환 실패"})
		return
	}
	user := tmp.(*dao.User)
	googleToken := oauth2.Token{
		AccessToken:  user.AccessToken,
		TokenType:    "Bearer",
		RefreshToken: user.RefreshToken,
	}
	contextTimeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	source := Config.TokenSource(contextTimeout, &googleToken)
	youtubeService, err := youtube.NewService(contextTimeout, option.WithTokenSource(source))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: "유튜브 서비스 호출 실패"})
		return
	}
	result, err := youtubeService.Channels.List(
		[]string{"auditDetails", "statistics", "snippet"},
	).Mine(true).Do()
	if len(result.Items) != 1 {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: "유튜브 서비스 조회 실패"})
		return
	}
	err = youtuber.ValidateChannel(result.Items[0])
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, Response{Message: "인증 실패사유: " + err.Error()})
		return
	}
	user.IsAuthedStreamer = true
	err = user.Save()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: "db업데이트 실패"})
		return
	}
	ctx.JSON(http.StatusOK, Response{Message: "성공"})
}
