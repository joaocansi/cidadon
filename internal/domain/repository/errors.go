package repository

import (
	"cidadon/internal/errors"
	stderrors "errors"
)

var (
	ErrDBNotFound = stderrors.New("resource not found")
	ErrDBConflict = stderrors.New("resource conflict")
	ErrDBInternal = stderrors.New("internal error")
)

var classification = map[error]errors.Code{
	ErrDBNotFound: errors.CodeNotFound,
	ErrDBConflict: errors.CodeConflict,
	ErrDBInternal: errors.CodeInternal,
}

func ClassifyCode(err error) errors.Code {
	for sentinel, code := range classification {
		if stderrors.Is(err, sentinel) {
			return code
		}
	}
	return errors.CodeInternal
}

func WrapDB(err error, message string) *errors.Error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, ClassifyCode(err), message)
}
