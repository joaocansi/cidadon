package service

import (
	stderrors "errors"
	"fmt"
)

type Code string

const (
	CodeInvalidInput    Code = "INVALID_INPUT"
	CodeUnauthorized    Code = "UNAUTHORIZED"
	CodeForbidden       Code = "FORBIDDEN"
	CodeNotFound        Code = "NOT_FOUND"
	CodeAlreadyExists   Code = "ALREADY_EXISTS"
	CodeConflict        Code = "CONFLICT"
	CodeFailedPrecond   Code = "FAILED_PRECOND"
	CodeResourceExhaust Code = "RESOURCE_EXHAUST"
	CodeTimeout         Code = "TIMEOUT"
	CodeUnavailable     Code = "UNAVAILABLE"
	CodeInternal        Code = "INTERNAL"
)

type Error struct {
	Code    Code
	Message string
	Details any
	Err     error
}

func (e *Error) Error() string {
	msg := e.Message
	if msg == "" && e.Err != nil {
		msg = e.Err.Error()
	}
	return msg
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) WithDetails(details any) *Error {
	e.Details = details
	return e
}

func New(code Code, message string) *Error {
	return &Error{Code: code, Message: message}
}

func Newf(code Code, format string, args ...any) *Error {
	return &Error{Code: code, Message: fmt.Sprintf(format, args...)}
}

func Wrap(err error, code Code, message string) *Error {
	if err == nil {
		return nil
	}
	return &Error{Code: code, Message: message, Err: err}
}

func InvalidInput(format string, args ...any) *Error {
	return Newf(CodeInvalidInput, format, args...)
}

func Unauthorized(format string, args ...any) *Error {
	return Newf(CodeUnauthorized, format, args...)
}

func Forbidden(format string, args ...any) *Error {
	return Newf(CodeForbidden, format, args...)
}

func NotFound(format string, args ...any) *Error {
	return Newf(CodeNotFound, format, args...)
}

func AlreadyExists(format string, args ...any) *Error {
	return Newf(CodeAlreadyExists, format, args...)
}

func Conflict(format string, args ...any) *Error {
	return Newf(CodeConflict, format, args...)
}

func Unavailable(format string, args ...any) *Error {
	return Newf(CodeUnavailable, format, args...)
}

func Internal(err error) *Error {
	return &Error{Code: CodeInternal, Message: "internal server error", Err: err}
}

func From(err error) (*Error, bool) {
	var appErr *Error
	if stderrors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

func GetCode(err error) Code {
	if appErr, ok := From(err); ok {
		return appErr.Code
	}
	return CodeInternal
}

func Is(err error, code Code) bool {
	return GetCode(err) == code
}
