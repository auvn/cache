package server

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrTimeout    Error = NewError("timeout", http.StatusGatewayTimeout)
	ErrBadRequest Error = NewError("bad request", http.StatusBadRequest)
	ErrInternal   Error = NewError("internal server error", http.StatusInternalServerError)
	ErrForbidden  Error = NewError("forbidden", http.StatusForbidden)
)

type Error interface {
	error
	Code() int
	Cause() error
	Internal() bool
	WithCause(cause error) Error
}

type serverError struct {
	err   error
	cause error
	code  int
}

func (self *serverError) copy() *serverError {
	return &serverError{
		err:   self.err,
		cause: self.cause,
		code:  self.code,
	}
}

func (self *serverError) Error() string {
	return fmt.Sprintf("%s: %s", self.err, self.cause)
}

func (self *serverError) Cause() error {
	return self.cause
}

func (self *serverError) Code() int {
	return self.code
}

func (self *serverError) Internal() bool {
	return self.code == http.StatusInternalServerError
}

func (self *serverError) WithCause(cause error) Error {
	err := self.copy()
	err.cause = cause
	return err
}

func NewError(err string, code int) Error {
	return &serverError{
		err:  errors.New(err),
		code: code,
	}
}
