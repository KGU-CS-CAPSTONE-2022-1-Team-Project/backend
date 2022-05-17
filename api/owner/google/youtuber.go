package google

import (
	"backend/infrastructure/owner/dao"
	"backend/internal/owner/youtuber"
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
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ResponseCommon{Message: "유저정보교환 실패"})
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
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ResponseCommon{
			"유튜브 서비스 호출 실패",
		})
		return
	}
	result, err := youtubeService.Channels.List(
		[]string{"auditDetails", "statistics", "snippet"},
	).Mine(true).Do()
	if len(result.Items) != 1 {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ResponseCommon{Message: "유튜브 서비스 조회 실패"})
		return
	}
	err = youtuber.ValidateChannel(result.Items[0])
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, ResponseCommon{Message: "인증 실패사유: " + err.Error()})
		return
	}
	user.IsAuthedStreamer = true
	err = user.Save()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ResponseCommon{Message: "db업데이트 실패"})
		return
	}
	ctx.Set("youtube", result.Items[0])
}

func SetChannel(ctx *gin.Context) {
	tmp, _ := ctx.Get("youtube")
	channel := tmp.(*youtube.Channel)
	tmp, _ = ctx.Get("user")
	user := tmp.(*dao.User)
	user.Channel.Name = channel.Snippet.Title
	user.Channel.Description = channel.Snippet.Description
	user.Channel.Url = "https://youtube.com/channel/" + channel.Id
	user.Channel.Image = channel.Snippet.Thumbnails.Default.Url
	if err := user.Save(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, ResponseCommon{
			Message: "유튜브 권한 부족",
		})
		return
	}
	ctx.JSON(http.StatusOK, ResponseCommon{Message: "성공"})
}

func GetChannel(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ResponseCommon{Message: "잘못된 파라미터"})
	}
	userDB := dao.User{ID: id}
	result, err := userDB.Read()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, ResponseCommon{Message: "존재하지 않는 Id"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"title":       result.Channel.Name,
		"description": result.Channel.Description,
		"image":       result.Channel.Image,
		"url":         result.Channel.Url,
	})
}
