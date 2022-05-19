package dao

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func IsEmpty(err error) bool {
	return errors.Cause(err) == mongo.ErrNoDocuments
}
