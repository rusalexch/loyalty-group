package app

import "errors"

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrRewardNotFound     = errors.New("reward not found")
	ErrProductNotFound    = errors.New("product not found")
	ErrRewardAlreadyExist = errors.New("reward already exist")
)
