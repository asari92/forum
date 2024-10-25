package entities

import (
	"errors"
)

var (
	ErrNoRecord = errors.New("no matching record found")

	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrDuplicateEmail = errors.New("duplicate email")

	ErrInvalidUser       = errors.New("invalid user")
	ErrInvalidToken      = errors.New("invalid session token")
	ErrInvalidCSRFToken  = errors.New("invalid csrf token")
)
