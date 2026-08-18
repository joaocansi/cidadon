package errors

import (
	stderrors "errors"
)

var (
	ErrDBNotFound = stderrors.New("resource not found")
	ErrDBConflict = stderrors.New("resource conflict")
	ErrDBInternal = stderrors.New("internal error")
)

var classification = map[error]Code{
	ErrDBNotFound: CodeNotFound,
	ErrDBConflict: CodeConflict,
	ErrDBInternal: CodeInternal,
}

func ClassifyCode(err error) Code {
	for sentinel, code := range classification {
		if stderrors.Is(err, sentinel) {
			return code
		}
	}
	return CodeInternal
}

func WrapDB(err error, message string) *Error {
	if err == nil {
		return nil
	}
	return Wrap(err, ClassifyCode(err), message)
}
