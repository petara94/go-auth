package api

import (
	"errors"
	"fmt"
)

type RestError struct {
	Message string `json:"message"`
}

func (r *RestError) Error() string {
	return r.Message
}

func RestErrorFromError(err error) error {
	return &RestError{Message: err.Error()}
}

func RestErrorf(str string, args ...any) error {
	return &RestError{Message: fmt.Sprintf(str, args...)}
}

var (
	ErrBadParam = errors.New("bad param")
	Success     = errors.New("success")
)
