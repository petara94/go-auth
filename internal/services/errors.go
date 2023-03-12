package services

import "errors"

var (
	ErrSessionExpired = errors.New("session expired")
	ErrLoginErr       = errors.New("wrong password or login")
)
