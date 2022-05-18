package google

import (
	"backend/infrastructure/owner/dao"
	"backend/internal/owner"
	"backend/proto/owner/pb"
	"backend/tool"
	"context"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	userProfile "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
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

type OwnerService struct {
	pb.UnimplementedOwnerServer
}

func (receiver *OwnerService) isValidate(tokenString string) error {
	accessToken := owner.AccessToken{}
	err := owner.Validate(&accessToken, tokenString)
	if err != nil {
		return errors.Wrap(err, "not validate")
	}
	return nil
}

func (receiver *OwnerService) exchangeGoogle(ctx context.Context, code string) (*oauth2.Token, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFunc()
	return Config.Exchange(timeout, code)
}

func (receiver *OwnerService) getGoogleEmail(ctx context.Context, token *oauth2.Token) (string, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, 10*time.Second)
	defer cancelFunc()
	source := Config.TokenSource(timeout, token)
	client, err := userProfile.NewService(timeout, option.WithTokenSource(source))
	userInfo, err := client.Userinfo.Get().Do()
	return userInfo.Email, err
}

func (receiver *OwnerService) Google(_ context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	url := Config.AuthCodeURL(
		req.Ip,
		oauth2.AccessTypeOffline,
	)
	return &pb.LoginResponse{
		AuthUrl: url,
	}, nil
}

func (receiver *OwnerService) GoogleCallBack(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	token, err := receiver.exchangeGoogle(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	email, err := receiver.getGoogleEmail(ctx, token)
	if err != nil {
		return nil, err
	}
	userDB := dao.User{Email: email}
	result, err := userDB.Read()
	userDB = dao.User{
		ID:           result.ID,
		Email:        email,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	if err == nil {
		err = userDB.Save()
		if err != nil {
			return nil, err
		}
	} else if dao.IsEmpty(err) {
		err = userDB.Create()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}
	accessToken := owner.AccessToken{
		UserID: userDB.ID,
	}
	tokenString, err := owner.CreateTokenString(&accessToken)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterResponse{
		AccessToken: tokenString,
	}, nil
}

func (receiver *OwnerService) SaveAddress(_ context.Context, req *pb.AddressRequest) (*pb.AddressResponse, error) {
	err := receiver.isValidate(req.AccessToken)
	if err != nil {
		return nil, err
	}
	accessToken := owner.AccessToken{}
	err = owner.GetAuthInfo(&accessToken, req.AccessToken)
	if err != nil {
		return nil, err
	}
	userDB := dao.User{ID: accessToken.UserID}
	result, err := userDB.Read()
	if err != nil {
		return nil, err
	}
	channel, err := GetYoutubeChannel(result)
	if err != nil {
		return nil, err
	}
	result.Address = req.Address
	SetChannel(result, channel)
	err = result.Save()
	if err != nil {
		return nil, err
	}
	go RegisterContract(req.Address, userDB.ID)
	return &pb.AddressResponse{
		IsValidate: true,
	}, nil
}

func (receiver OwnerService) GetChannel(_ context.Context, req *pb.ChannelRequest) (*pb.ChannelResponse, error) {
	userDB := dao.User{ID: req.Id}
	result, err := userDB.Read()
	if err != nil {
		return nil, err
	}
	return &pb.ChannelResponse{
		Title:       result.Channel.Name,
		Description: result.Channel.Description,
		Image:       result.Channel.Image,
		Url:         result.Channel.Url,
	}, nil
}
