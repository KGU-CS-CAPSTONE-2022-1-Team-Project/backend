package owner

import (
	"backend/infrastructure/owner/dao"
)

type User struct {
	ID               string
	Email            string
	IsAuthedStreamer bool
	AccessToken      string
	RefreshToken     string
	Address          string
	Channel
}

type Channel struct {
	Name        string
	Description string
	image       string
	URL         string
}

func User2UserDB(user User) dao.User {
	return dao.User{
		ID:               user.ID,
		Email:            user.Email,
		IsAuthedStreamer: user.IsAuthedStreamer,
		AccessToken:      user.AccessToken,
		RefreshToken:     user.RefreshToken,
	}
}
