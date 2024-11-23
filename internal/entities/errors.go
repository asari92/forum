package entities

import (
	"errors"
)

var (
	ErrNoRecord = errors.New("no matching record found")

	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidData        = errors.New("invalid data")

	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUsername = errors.New("duplicate username")

	ErrInvalidUser      = errors.New("invalid user")
	ErrInvalidToken     = errors.New("invalid session token")
	ErrInvalidCSRFToken = errors.New("invalid csrf token")
)
