package serializer

import (
	"errors"
	"fmt"
)

type Error interface {
	error
	Cause() error
	WithCause(cause error) Error
}

type sError struct {
	err   error
	cause error
}

func (self *sError) copy() *sError {
	return &sError{err: self.err, cause: self.err}
}

func (self *sError) Cause() error {
	return self.cause
}

func (self *sError) WithCause(cause error) Error {
	serr := self.copy()
	serr.cause = cause
	return serr
}

func (self *sError) Error() string {
	if self.cause == nil {
		return fmt.Sprintf("%s", self.err.Error())
	} else {
		return fmt.Sprintf("%s: %s", self.err.Error(), self.cause.Error())
	}
}

func NewError(err string) Error {
	return &sError{err: errors.New(err)}
}

type prefixError struct {
	err Error
}

func (self *prefixError) copy() *prefixError {
	return &prefixError{
		err: self.err,
	}
}

func (self *prefixError) Cause() error {
	return self.err.Cause()
}

func (self *prefixError) WithCause(cause error) Error {
	return self.err.WithCause(cause)
}

func (self *prefixError) Error() string {
	return self.err.Error()
}

func (self *prefixError) WithGot(r rune) *prefixError {
	err := self.copy()
	err.err = err.err.WithCause(fmt.Errorf("got %q", r))
	return err
}

func newPrefixError(expected rune) *prefixError {
	return &prefixError{
		err: NewError(fmt.Sprintf("expected prefix %q", expected)),
	}
}
