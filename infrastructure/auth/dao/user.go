package dao

import (
	"errors"
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
	Address          string
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

func (receiver *User) Save() error {
	db, err := dbConnection()
	if err != nil {
		return err
	}
	return db.Save(receiver).Error
}

func (receiver *User) Read() (*User, error) {
	db, err := dbConnection()
	if err != nil {
		return nil, err
	}
	if receiver.ID != "" {
		if receiver.TokenIdentifier != "" {
			return receiver.readByTokenID(db)
		} else {
			return receiver.readByUserID(db)
		}
	}
	if receiver.Email != "" {
		return receiver.readByEmail(db)
	}
	return nil, errors.New("찾을 수 없는 조건")
}

func (receiver *User) readByUserID(conn *gorm.DB) (*User, error) {
	result := &User{}
	return result, conn.First(&result, "id", receiver.ID).Error
}

func (receiver *User) readByTokenID(conn *gorm.DB) (*User, error) {
	result := &User{}
	return result, conn.First(&result, "id=? AND token_identifier=?",
		receiver.ID, receiver.TokenIdentifier).Error
}

func (receiver *User) readByEmail(conn *gorm.DB) (*User, error) {
	result := &User{}
	return result, conn.First(&result, "email=?",
		receiver.Email).Error
}
