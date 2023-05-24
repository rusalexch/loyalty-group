package auth

import "errors"

var (
	errLoginAlreadyExist = errors.New("login already exist")
	errUnauthorized      = errors.New("unauthorized")
)
