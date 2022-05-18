package dao

import (
	"errors"
	"github.com/google/uuid"
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
		panic("Fail migration")
	}
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
		return receiver.readByUserID(db)
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

func (receiver *User) readByEmail(conn *gorm.DB) (*User, error) {
	result := &User{}
	return result, conn.First(&result, "email=?",
		receiver.Email).Error
}
