package dao

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type User struct {
	ID               string `gorm:"primaryKey'"`
	Email            string `gorm:"not_null"`
	IsAuthedStreamer bool   `gorm:"not_null"`
	AccessToken      string
	RefreshToken     string
	Address          string
	Channel
	gorm.Model
}

type Channel struct {
	Name        string
	Description string
	Image       string
	Url         string
}

func init() {
	user := User{}
	err := user.Migration()
	if err != nil {
		panic(errors.Wrap(err, "init"))
	}
}
