package commands

import (
	"errors"

	"github.com/auvn/go.cache/session"
)

var (
	ErrAuthRequired = errors.New("auth required")

	ErrUnknownCommand    = errors.New("unknown command")
	ErrNumberOfArguments = errors.New("wrong number of arguments")
	ErrValueObjectIsNil  = errors.New("value object is nil")
	ErrNotFound          = errors.New("value not found")
	ErrWrongType         = errors.New("accessing a key holding the wrong type of value")
	ErrNonStr            = errors.New("non str")
	ErrNonInt            = errors.New("non int")
)

type Command interface {
	Execute(session.Session, Arguments) (interface{}, error)
	IsFlag(flag int) bool
	Flag() int
}

type Registry interface {
	Get(name string) (Command, bool)
}

type MapRegistry map[string]Command

func (self MapRegistry) Get(name string) (Command, bool) {
	cmd, ok := self[name]
	return cmd, ok
}
