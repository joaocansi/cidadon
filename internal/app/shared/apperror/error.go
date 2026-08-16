package apperror

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ErrorType string

const (
	NotFound     ErrorType = "not_found"
	Conflict     ErrorType = "conflict"
	Unauthorized ErrorType = "unauthorized"
	Validation   ErrorType = "validation"
	Internal     ErrorType = "internal"
)

type Error struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Details any       `json:"details,omitempty"`
	Err     error     `json:"-"`
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) StatusCode() int {
	switch e.Type {
	case Validation:
		return 400
	case Unauthorized:
		return 401
	case NotFound:
		return 404
	case Conflict:
		return 409
	default:
		return 500
	}
}

const (
	CodeRefreshTokenNotFound = "REFRESH_TOKEN_NOT_FOUND"
	CodeAccessTokenExpired   = "ACCESS_TOKEN_EXPIRED"
	CodeAccessTokenInvalid   = "ACCESS_TOKEN_INVALID"
)

func NewUnauthorizedWithCode(message, code string) *Error {
	return &Error{
		Type:    Unauthorized,
		Message: message,
		Details: map[string]interface{}{"code": code},
	}
}

func NewNotFound(message string) *Error {
	return &Error{
		Type:    NotFound,
		Message: message,
	}
}

func NewConflict(message string) *Error {
	return &Error{
		Type:    Conflict,
		Message: message,
	}
}

func NewUnauthorized(message string) *Error {
	return &Error{
		Type:    Unauthorized,
		Message: message,
	}
}

func NewValidation(message string, details any) *Error {
	return &Error{
		Type:    Validation,
		Message: message,
		Details: details,
	}
}

func NewInternal(message string, err ...error) *Error {
	return &Error{
		Type:    Internal,
		Message: message,
		Err:     fmt.Errorf(message),
	}
}

func FromError(err error) *Error {
	if err == nil {
		return nil
	}

	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr
	}

	if validationErr := fromValidation(err); validationErr != nil {
		return validationErr
	}

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return NewNotFound("resource not found")

	case errors.Is(err, gorm.ErrDuplicatedKey):
		return NewConflict("resource already exists")
	}

	return NewInternal("internal server error", err)
}

type ValidationDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func fromValidation(err error) *Error {
	var ve validator.ValidationErrors

	if !errors.As(err, &ve) {
		return nil
	}

	details := make([]ValidationDetail, 0, len(ve))

	for _, fe := range ve {
		details = append(details, ValidationDetail{
			Field:   fe.Field(),
			Message: validationMessage(fe),
		})
	}

	return NewValidation("validation failed", details)
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "min":
		return fmt.Sprintf("must have at least %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("must have at most %s characters", fe.Param())
	case "strong_password":
		return "must contain an uppercase letter, lowercase letter, number, and special character"
	case "transaction_type":
		return "must be one of: income, expense, refund, transfer"
	case "business_day_rule":
		return "must be one of: none, next or previous"
	case "gte":
		return fmt.Sprintf("must be greater than or equal to %s", fe.Param())
	case "lte":
		return fmt.Sprintf("must be less than or equal to %s", fe.Param())
	default:
		return "is invalid"
	}
}
