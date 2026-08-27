package repository

type DBErrorCode uint8

const (
	DBErrorNotFound DBErrorCode = iota
	DBErrorConflict
	DBErrorInternal
)

type DBError struct {
	Code DBErrorCode
	Err  error
}

func NewDBError(code DBErrorCode, err error) *DBError {
	return &DBError{
		Code: code,
		Err:  err,
	}
}

func (e *DBError) Error() string {
	return e.Err.Error()
}

func (e *DBError) Unwrap() error {
	return e.Err
}
