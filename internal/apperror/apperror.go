package apperror

import "encoding/json"

var (
	ErrDataBaseConnect    = NewAppError(nil, "Connection error", "DB-000001")
	ErrDataNotFound       = NewAppError(nil, "Not found", "MDL-000001")
	ErrConfigNotFoundFile = NewAppError(nil, "Not found config file", "CFG-000001")
	ErrCacheKeyNotFound   = NewAppError(nil, "Key not found in cache", "CCH-000001")
)

type AppError struct {
	Err     error  `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewAppError(err error, message, code string) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}
