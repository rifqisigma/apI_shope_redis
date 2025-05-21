package helper

import "errors"

var (
	ErrInternal = errors.New("internal error")

	//auth
	ErrInvalidEmail = errors.New("email tidak sesuai")

	//shop
	ErrNotAdmin       = errors.New("kau bukan admin")
	ErrStocknotEnough = errors.New("stock tidak cukup")
)
