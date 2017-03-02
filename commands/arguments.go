package commands

import (
	"github.com/auvn/go.cache/core"
)

type Arguments interface {
	Len() int
	Iter() ArgumentsIterator
	IterN(n int) (ArgumentsIterator, error)
	IterAtLeast(n int) (ArgumentsIterator, error)
}

type ArgumentsIterator interface {
	Next() (core.Value, error)
	NextStr() (core.StrValue, error)
	NextInt() (core.IntValue, error)
	NextArray() core.ValueArray
	NextArguments() Arguments
}

type argsIterator struct {
	cur  int
	args []core.Value
	len  int
}

func (self *argsIterator) Next() (core.Value, error) {
	if self.cur >= self.len {
		return core.Value{}, ErrNumberOfArguments
	}
	pos := self.cur
	self.cur += 1
	return self.args[pos], nil
}

func (self *argsIterator) NextStr() (core.StrValue, error) {
	str := core.EmptyStrValue
	val, err := self.Next()
	if err != nil {
		return str, err
	}

	str, err = val.Str()
	if err != nil {
		return str, ErrNonStr
	}

	return str, nil
}

func (self *argsIterator) NextInt() (core.IntValue, error) {
	i := core.EmptyIntValue
	val, err := self.Next()
	if err != nil {
		return i, err
	}

	i, err = val.Int()
	if err != nil {
		return i, ErrNonInt
	}

	return i, nil
}

func (self *argsIterator) NextArray() core.ValueArray {
	arr := make(core.ValueArray, 0, self.len)
	stop := false
	for !stop {
		if v, err := self.Next(); err == nil {
			arr = append(arr, v)
		} else {
			stop = true
		}
	}
	return arr
}

func (self *argsIterator) NextArguments() Arguments {
	values := self.NextArray()
	return NewArguments(values...)
}

type arguments []core.Value

func (self arguments) Len() int {
	return len(self)
}

func (self arguments) Iter() ArgumentsIterator {
	iter, _ := self.IterN(self.Len())
	return iter
}

func (self arguments) IterN(n int) (ArgumentsIterator, error) {
	len := self.Len()
	if n != len {
		return nil, ErrNumberOfArguments
	}
	return &argsIterator{args: self, len: len}, nil
}

func (self arguments) IterAtLeast(n int) (ArgumentsIterator, error) {
	len := self.Len()
	if n > len {
		return nil, ErrNumberOfArguments
	}
	return &argsIterator{args: self, len: len}, nil
}

func NewArguments(values ...core.Value) Arguments {
	return append(make(arguments, 0, len(values)), values...)
}
