package database

import (
	"errors"
)

var (
	ErrDBNotFound = errors.New("resource not found")
	ErrDBConflict = errors.New("resource conflict")
	ErrDBInternal = errors.New("internal error")
)
