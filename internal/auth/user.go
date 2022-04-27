package auth

import "backend/internal/auth/dao"

type User struct {
	ID               string
	Email            string
	TokenIdentifier  string
	IsAuthedStreamer bool
	AccessToken      string
	RefreshToken     string
}

func UserDB2User(user dao.User) User {
	return User{
		ID:               user.ID,
		Email:            user.Email,
		TokenIdentifier:  user.TokenIdentifier,
		IsAuthedStreamer: user.IsAuthedStreamer,
	}
}

func User2UserDB(user User) dao.User {
	return dao.User{
		ID:               user.ID,
		Email:            user.Email,
		TokenIdentifier:  user.TokenIdentifier,
		IsAuthedStreamer: user.IsAuthedStreamer,
		AccessToken:      user.AccessToken,
		RefreshToken:     user.RefreshToken,
	}
}
