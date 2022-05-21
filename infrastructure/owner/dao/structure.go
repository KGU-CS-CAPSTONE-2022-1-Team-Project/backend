package dao

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type Owner struct {
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

type User struct {
	Address   string `gorm:"primaryKey"`
	Nickname  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func init() {
	owner := Owner{}
	user := User{}
	err := owner.Migration()
	if err != nil {
		panic(errors.Wrap(err, "init"))
	}
	err = user.Migration()
	if err != nil {
		panic(errors.Wrap(err, "init"))
	}
}
