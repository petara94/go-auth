package services

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"
)

var (
	ErrSessionExpired     = errors.New("session expired")
	ErrLoginErr           = errors.New("wrong password or login")
	ErrLoginAlreadyExists = errors.New("user with same login already exists")
	ErrWrongPassword      = errors.New("wrong password")
	ErrWrongPagination    = errors.New("wrong pagination")
	ErrWeakPassword       = errors.New("weak password")
	ErrSamePassword       = errors.New("same password as old one")
)

// handle pgx constraint errors
func HandleRepositoryError(err error) error {
	var pgErr *pgconn.PgError

	if !errors.As(err, &pgErr) {
		return err
	}

	// check constraint violation
	if !isConstraintViolation(pgErr.Code) {
		return err
	}

	switch pgErr.ConstraintName {
	case "users_login_ukey":
		return ErrLoginAlreadyExists
	default:
		return err
	}
}

// check constraint violation func
func isConstraintViolation(code string) bool {
	// see https://www.postgresql.org/docs/11/errcodes-appendix.html
	const constraintCodePrefix = "23"

	return strings.HasPrefix(code, constraintCodePrefix)
}
