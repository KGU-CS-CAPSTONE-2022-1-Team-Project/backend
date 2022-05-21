package dao

import (
	"github.com/pkg/errors"
)

func (receiver *User) Migration() error {
	db, err := dbConnection()
	if err != nil {
		return errors.Wrap(err, "Migration")
	}
	err = db.Migrator().AutoMigrate(receiver)
	if err != nil {
		return err
	}
	return nil
}

func (receiver *User) Load() error {
	db, err := dbConnection()
	if err != nil {
		return nil
	}
	if receiver.Address == "" && receiver.Nickname == "" {
		return errors.New("invalidate param")
	}
	return db.First(receiver).Error
}

func (receiver *User) Save() error {
	db, err := dbConnection()
	if err != nil {
		return errors.Wrap(err, "User.Save")
	}
	return db.Create(receiver).Error
}

func (receiver *User) Read() error {
	db, err := dbConnection()
	if err != nil {
		return errors.Wrap(err, "User.Read")
	}
	err = db.First(receiver).Error
	if err != nil {
		return errors.Wrap(err, "User.Read")
	}
	return nil
}
