package serializer

import "strconv"

var (
	ErrPayloadNonArray  = NewError("not-array payload received")
	ErrPayloadNonValue  = NewError("not-value payload received")
	ErrPayloadNonString = NewError("not-string payload received")
	ErrPayloadNonInt    = NewError("not-int payload received")
	ErrPayloadNonBool   = NewError("not-bool payload received")
	ErrPayloadNonError  = NewError("not-error payload received")
)

type Payload interface {
	Array() ([]Payload, error)
	Bytes() ([]byte, error)
	Str() (string, error)
	Int() (int, error)
	Bool() (bool, error)
	Err() error

	IsNil() bool
	IsArray() bool
	IsErr() bool
}

type payload struct {
	v interface{}
}

func (self *payload) Array() ([]Payload, error) {
	if arr, ok := self.v.([]Payload); ok {
		return arr, nil
	}
	return []Payload{}, ErrPayloadNonArray
}

func (self *payload) Bytes() ([]byte, error) {
	if bs, ok := self.v.([]byte); ok {
		return bs, nil
	}
	return []byte{}, ErrPayloadNonValue
}

func (self *payload) Str() (string, error) {
	if bs, err := self.Bytes(); err == nil {
		return string(bs), nil
	}
	return "", ErrPayloadNonString
}

func (self *payload) Int() (int, error) {
	if s, err := self.Str(); err == nil {
		if i, err := strconv.Atoi(string(s)); err == nil {
			return i, nil
		}
	}
	return 0, ErrPayloadNonInt
}

func (self *payload) Bool() (bool, error) {
	if bs, err := self.Bytes(); err == nil && len(bs) == 1 {
		return bs[0] != '0', nil
	}
	return false, ErrPayloadNonBool
}

func (self *payload) Err() error {
	if err, ok := self.v.(error); ok {
		return err
	}
	return ErrPayloadNonError
}

func (self *payload) IsNil() bool {
	return self.v == nil
}

func (self *payload) IsArray() bool {
	_, ok := self.v.([]Payload)
	return ok
}

func (self *payload) IsErr() bool {
	_, ok := self.v.(error)
	return ok
}
