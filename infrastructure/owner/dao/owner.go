package dao

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// BeforeCreate is gorm Hook. plz not call.
func (receiver *Original) BeforeCreate(_ *gorm.DB) (err error) {
	receiver.ID = uuid.NewString()
	return nil
}

func (receiver *Original) Migration() error {
	err := db.Migrator().AutoMigrate(receiver)
	if err != nil {
		return err
	}
	return nil
}

func (receiver *Original) Create() error {
	err := db.Create(receiver).Error
	if err != nil {
		return err
	}
	return nil
}

func (receiver *Original) Save() error {
	return db.Model(&Original{ID: receiver.ID}).
		Updates(&receiver).Error
}

func (receiver *Original) Read() (*Original, error) {
	if receiver.ID != "" {
		return receiver.readByUserID(db)
	}
	if receiver.Email != "" {
		return receiver.readByEmail(db)
	}
	if receiver.Address != "" {
		return receiver.readByChannelName(db)
	}
	return nil, errors.New("찾을 수 없는 조건")
}

func (receiver *Original) readByUserID(conn *gorm.DB) (*Original, error) {
	result := &Original{}
	return result, conn.First(&result, "id", receiver.ID).Error
}

func (receiver *Original) readByEmail(conn *gorm.DB) (*Original, error) {
	result := &Original{}
	return result, conn.First(&result, "email=?",
		receiver.Email).Error
}

func (receiver *Original) readByChannelName(conn *gorm.DB) (*Original, error) {
	result := &Original{}
	return result, conn.First(&result, "address=?",
		receiver.Email).Error
}
