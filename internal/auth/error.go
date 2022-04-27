package auth

import "errors"

var (
	ErrValidateToken = errors.New("invalidate token")
	ErrDecrypt       = errors.New("fail decrypt")
)
