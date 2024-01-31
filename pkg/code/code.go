package code

import (
	"errors"
)

// constants
const (
	OK                   = 0
	ParamIncorrect       = 1000
	NotFound             = 1001
	Timeout              = 1002
	InternalUnknownError = 2999
)

// define errors
var (
	ErrParamIncorrect = errors.New("param incorrect")
	ErrInternalError  = errors.New("internal error")
)

// CustomError a custom error
type CustomError struct {
	Code       int   `json:"code"`
	Error      error `json:"error"`
	HttpStatus int   `json:"http_status"`
}

// NewCustomError create a new CustomError
func NewCustomError(code, httpStatus int, err error) *CustomError {
	return &CustomError{
		Code:       code,
		Error:      err,
		HttpStatus: httpStatus,
	}
}
