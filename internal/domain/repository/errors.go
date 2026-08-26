package repository

import (
	"cidadon/internal/domain/service"
	stderrors "errors"
)

var (
	ErrDBNotFound = stderrors.New("resource not found")
	ErrDBConflict = stderrors.New("resource conflict")
	ErrDBInternal = stderrors.New("internal error")
)

func ClassifyCode(err error) service.Code {
	switch {
	case stderrors.Is(err, ErrDBNotFound):
		return service.CodeNotFound
	case stderrors.Is(err, ErrDBConflict):
		return service.CodeConflict
	case stderrors.Is(err, ErrDBInternal):
		return service.CodeInternal
	default:
		return service.CodeInternal
	}
}

func WrapDB(err error, message string) *service.Error {
	if err == nil {
		return nil
	}
	return service.Wrap(err, ClassifyCode(err), message)
}
