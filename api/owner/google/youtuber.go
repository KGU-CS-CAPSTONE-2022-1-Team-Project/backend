package google

import (
	"backend/infrastructure/owner/dao"
	"backend/internal/owner/youtuber"
	"backend/tool"
	"context"
	"errors"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"time"
)

var whiteList []string

func init() {
	if len(whiteList) == 0 {
		tool.ReadConfig("./configs/owner", "client_secret", "json")
	}
	whiteList = viper.GetStringSlice("white_list")
}

func GetYoutubeChannel(user *dao.Original) (*youtube.Channel, error) {
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
		return nil, err
	}
	result, err := youtubeService.Channels.List(
		[]string{"auditDetails", "statistics", "snippet"},
	).Mine(true).Do()
	if len(result.Items) != 1 {
		return nil, errors.New("not found")
	}
	isWhiteList := false
	for _, white := range whiteList {
		if user.Email == white {
			isWhiteList = true
			break
		}
	}
	if !isWhiteList {
		err = youtuber.ValidateChannel(result.Items[0])
		if err != nil {
			tool.Logger().Warning("invalidate channel", err, "uid", user.ID)
			return nil, err
		}
	}
	tool.Logger().Info("validate channel", "uid", user.ID)
	return result.Items[0], nil
}

func SetChannel(user *dao.Original, channel *youtube.Channel) {
	user.Channel.Name = channel.Snippet.Title
	user.Channel.Description = channel.Snippet.Description
	user.Channel.Url = "https://youtube.com/channel/" + channel.Id
	user.Channel.Image = channel.Snippet.Thumbnails.Default.Url
}
