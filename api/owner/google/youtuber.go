package google

import (
	"backend/infrastructure/owner/dao"
	"context"
	"errors"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"time"
)

func GetYoutubeChannel(user *dao.Owner) (*youtube.Channel, error) {
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
		return nil, errors.New("Not Found")
	}
	//err = youtuber.ValidateChannel(result.Items[0])
	//if err != nil {
	//	return nil, err
	//}
	return result.Items[0], nil
}

func SetChannel(user *dao.Owner, channel *youtube.Channel) {
	user.Channel.Name = channel.Snippet.Title
	user.Channel.Description = channel.Snippet.Description
	user.Channel.Url = "https://youtube.com/channel/" + channel.Id
	user.Channel.Image = channel.Snippet.Thumbnails.Default.Url
}
