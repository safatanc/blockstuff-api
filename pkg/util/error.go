package util

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type Error struct {
	StatusCode int
	Err        error
}

func NewError(statusCode int, err error) *Error {
	return &Error{
		StatusCode: statusCode,
		Err:        err,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("status %v: %v", e.StatusCode, e.Err)
}

func GetErrorStatusCode(err error) int {
	defaultStatusCode := 500

	// Postgres Error
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return http.StatusConflict
	}

	// GORM Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound
	} else if errors.Is(err, gorm.ErrDuplicatedKey) {
		return http.StatusConflict
	}

	// Validation Error
	if strings.Contains(err.Error(), "validation") {
		return http.StatusBadRequest
	}

	return defaultStatusCode
}
