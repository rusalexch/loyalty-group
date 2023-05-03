package common

import "errors"

var (
	ErrProductNotFound     = errors.New("product not found")
	ErrProductAlreadyExist = errors.New("product already exist")
)
