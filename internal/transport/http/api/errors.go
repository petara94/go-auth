package api

import (
	"errors"
	"fmt"
)

type RestError struct {
	Message error `json:"message"`
}

func (r *RestError) Error() string {
	return r.Message.Error()
}

func RestErrorFromError(err error) error {
	return &RestError{Message: err}
}

func RestErrorf(str string, args ...any) error {
	return &RestError{Message: fmt.Errorf(str, args...)}
}

var (
	ErrBadParam = errors.New("bad param")
)
