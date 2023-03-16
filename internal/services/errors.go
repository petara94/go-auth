package services

import "errors"

var (
	ErrSessionExpired     = errors.New("session expired")
	ErrLoginErr           = errors.New("wrong password or login")
	ErrLoginAlreadyExists = errors.New("user with same login already exists")
)
