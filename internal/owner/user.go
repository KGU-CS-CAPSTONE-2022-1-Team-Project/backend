package owner

import (
	"backend/infrastructure/owner/dao"
	"encoding/hex"
	"github.com/pkg/errors"
	"strings"
	"unicode/utf8"
)

func (receiver *User) validate() error {
	if receiver.Nickname == "" || receiver.Address == "" {
		return errors.Wrap(errors.New("empty"), "Nickname Validator")
	}
	addr := strings.TrimPrefix(receiver.Address, "0x")
	addr = strings.TrimPrefix(addr, "0x")
	receiver.Nickname = strings.TrimPrefix(receiver.Nickname, " ")
	bytes, err := hex.DecodeString(addr)
	if err != nil {
		return errors.Wrap(errors.New("not hex address"), "Nickname Validator")
	}
	if len(bytes) != 20 {
		return errors.Wrap(errors.New("not klaytn addr"), "Nickname Validator")
	}
	if utf8.RuneCountInString(receiver.Nickname) < 4 || utf8.RuneCountInString(receiver.Nickname) > 10 {
		return errors.Wrap(errors.New("invalidate nickname len"), "Nickname Validator")
	}
	return nil
}

func User2UserDB(user User) dao.User {
	return dao.User{
		Address:  user.Address,
		Nickname: user.Nickname,
	}
}
