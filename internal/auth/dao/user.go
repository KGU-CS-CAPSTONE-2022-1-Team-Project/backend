package dao

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID               string `gorm:"primaryKey'"`
	Email            string `gorm:"not_null"`
	TokenIdentifier  string `gorm:"not_null"`
	IsAuthedStreamer bool   `gorm:"not_null"`
	AccessToken      string
	RefreshToken     string
	gorm.Model
}

// BeforeCreate is gorm Hook. plz not call.
func (receiver *User) BeforeCreate(_ *gorm.DB) (err error) {
	receiver.ID = uuid.NewString()
	return nil
}

func (receiver *User) Migration() error {
	db, err := dbConnection()
	if err != nil {
		return err
	}
	err = db.Migrator().AutoMigrate(receiver)
	if err != nil {
		return err
	}
	return nil
}

func (receiver *User) Create() error {
	db, err := dbConnection()
	if err != nil {
		return err
	}
	err = db.Create(receiver).Error
	if err != nil {
		return err
	}
	return nil
}

func (receiver *User) ReadByID() (*User, error) {
	db, err := dbConnection()
	if err != nil {
		return nil, err
	}
	result := &User{}
	err = db.First(result, "id", receiver.ID).Error
	return result, err
}

func (receiver User) ReadByTokenID(tokenId string) (*User, error) {
	db, err := dbConnection()
	if err != nil {
		return nil, err
	}
	result := &User{}
	err = db.First(result, "id=? AND token_identifier=?", receiver.ID, tokenId).Error
	return result, err
}
