package dao

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func IsEmpty(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
