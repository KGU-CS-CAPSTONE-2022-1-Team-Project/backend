package dao

import (
	"github.com/pkg/errors"
)

func (receiver *User) Migration() error {
	err := db.Migrator().AutoMigrate(receiver)
	if err != nil {
		return err
	}
	return nil
}

func (receiver *User) Load() error {
	if receiver.Address == "" && receiver.Nickname == "" {
		return errors.New("invalidate param")
	}
	return db.First(receiver).Error
}

func (receiver *User) Save() error {
	return db.Create(receiver).Error
}

func (receiver *User) Read() error {
	if receiver.Address != "" {
		if receiver.Nickname != "" {
			return errors.Wrap(receiver.readByAll(), "readByAll")
		}
		return errors.Wrap(receiver.readByAddr(), "readByAddr")
	} else if receiver.Nickname != "" {
		return errors.Wrap(receiver.readByNickname(), "readByNickname")
	}
	err := db.First(receiver).Error
	if err != nil {
		return errors.Wrap(err, "User.Read")
	}
	return nil
}

func (receiver *User) readByAddr() error {
	receiver.Nickname = ""
	return db.First(receiver).Error
}

func (receiver *User) readByNickname() error {
	return db.First(receiver, "nickname=?", receiver.Nickname).Error
}
func (receiver *User) readByAll() error {
	address := receiver.Address
	nickname := receiver.Nickname
	receiver.Nickname = ""
	receiver.Address = ""
	return db.First(receiver, "nickname=? OR address=?",
		nickname,
		address).Error
}
