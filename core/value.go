package core

import "strconv"

type Value []byte

func (self Value) String() string {
	return string(self)
}

func (self Value) Str() (StrValue, error) {
	return StrValue(self), nil
}

func (self Value) Int() (IntValue, error) {
	i, err := self.int(10, 0)
	return IntValue(i), err
}

func (self Value) Int64() (int64, error) {
	return self.int(10, 64)
}

func (self Value) int(base int, bitSize int) (int64, error) {
	return strconv.ParseInt(self.String(), base, bitSize)
}

func (self Value) Bytes() []byte {
	return []byte(self)
}

type IntValue int

func (self IntValue) Value() int {
	return int(self)
}

type StrValue string

func (self StrValue) Value() string {
	return string(self)
}

type ErrValue error

type ValueArray []Value

func (self ValueArray) Value() []Value {
	return []Value(self)
}

var (
	EmptyIntValue = IntValue(0)
	EmptyStrValue = StrValue("")
	EmptyValue    = Value{}
)
