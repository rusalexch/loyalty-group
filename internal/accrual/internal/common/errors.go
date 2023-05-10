package common

import "errors"

var (
	ErrOrderNotFound       = errors.New("order not found")
	ErrRewardNotFound      = errors.New("reward not found")
	ErrProductNotFound     = errors.New("product not found")
	ErrProductAlreadyExist = errors.New("product already exist")
)
