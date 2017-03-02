package client

import (
	"errors"

	"github.com/auvn/go.cache/net/serializer"
)

var (
	emptySlice       = []serializer.Payload{}
	emptyBytesSlice  = [][]byte{}
	emptyBytes       = []byte{}
	emptyStringSlice = []string{}
)

type MultiPayload []serializer.Payload

func (self MultiPayload) Array() ([]serializer.Payload, error) {
	var (
		arr []serializer.Payload
		err error
	)

	ret := make([]serializer.Payload, 0, 30)
	for _, p := range self {
		if arr, err = p.Array(); err != nil {
			return emptySlice, err
		} else {
			ret = append(ret, arr...)
		}
	}
	return ret, nil
}

func (self MultiPayload) Bytes() ([]byte, error) {
	if len(self) != 1 {
		return nil, errors.New("Bytes() cannot be invoked if len(multipayload) != 1")
	}
	return self[0].Bytes()
}

func (self MultiPayload) Str() (string, error) {
	if len(self) != 1 {
		return "", errors.New("Str() cannot be invoked if len(multipayload) != 1")
	}
	return self[0].Str()
}

func (self MultiPayload) Int() (int, error) {
	var (
		v   int
		err error
	)

	ret := 0
	for _, p := range self {
		if v, err = p.Int(); err != nil {
			return ret, err
		} else {
			ret += v
		}
	}
	return ret, nil
}

func (self MultiPayload) Bool() (bool, error) {
	var (
		b   bool
		err error
	)
	for _, p := range self {
		if b, err = p.Bool(); err != nil {
			return b, err
		} else if !b {
			return false, nil
		}
	}
	return true, nil
}

func (self MultiPayload) Err() error {
	var (
		err error
	)
	for _, p := range self {
		if err = p.Err(); err != nil && err != serializer.ErrPayloadNonError {
			return err
		}
	}
	return nil
}

func (self MultiPayload) IsNil() bool {
	for _, p := range self {
		if !p.IsNil() {
			return false
		}
	}
	return true
}

func (self MultiPayload) IsArray() bool {
	for _, p := range self {
		if !p.IsArray() {
			return false
		}
	}
	return true
}

func (self MultiPayload) IsErr() bool {
	for _, p := range self {
		if p.IsErr() {
			return true
		}
	}
	return false
}

type BoolCommand interface {
	Bool() (bool, error)
}

type IntCommand interface {
	Int() (int, error)
}

type BytesCommand interface {
	Bytes() ([]byte, error)
}

type BytesSliceCommand interface {
	BytesSlice() ([][]byte, error)
}

type StringSliceCommand interface {
	StringSlice() ([]string, error)
}

type Command interface {
	BoolCommand
	IntCommand
	BytesCommand
	BytesSliceCommand
	StringSliceCommand
}

type Payload []interface{}

type Caller interface {
	Call(cmdDef *CommandDefinition) (serializer.Payload, error)
}

type RemoteCommand struct {
	cmdDef *CommandDefinition
	caller Caller
}

func (self *RemoteCommand) call() (serializer.Payload, error) {
	return self.caller.Call(self.cmdDef)
}

func (self *RemoteCommand) Bool() (bool, error) {
	if res, err := self.call(); err != nil {
		return false, err
	} else {
		return res.Bool()
	}
}

func (self *RemoteCommand) Int() (int, error) {
	if res, err := self.call(); err != nil {
		return 0, err
	} else {
		return res.Int()
	}
}
func (self *RemoteCommand) Bytes() ([]byte, error) {
	if res, err := self.call(); err != nil {
		return emptyBytes, err
	} else if res.IsNil() {
		return emptyBytes, err
	} else {
		return res.Bytes()
	}
}

func (self *RemoteCommand) slice() ([]serializer.Payload, error) {
	if res, err := self.call(); err != nil {
		return emptySlice, err
	} else if res.IsNil() {
		return emptySlice, nil
	} else {
		if arr, err := res.Array(); err != nil {
			return emptySlice, err
		} else {
			return arr, nil
		}
	}
}
func (self *RemoteCommand) BytesSlice() ([][]byte, error) {
	arr, err := self.slice()
	if err != nil {
		return emptyBytesSlice, err
	}
	ret := make([][]byte, len(arr))
	for i, p := range arr {
		if bs, err := p.Bytes(); err != nil {
			return emptyBytesSlice, err
		} else {
			ret[i] = bs
		}
	}
	return emptyBytesSlice, nil
}

func (self *RemoteCommand) StringSlice() ([]string, error) {
	arr, err := self.slice()
	if err != nil {
		return emptyStringSlice, err
	}
	ret := make([]string, len(arr))
	for i, p := range arr {
		if s, err := p.Str(); err != nil {
			return emptyStringSlice, err
		} else {
			ret[i] = s
		}
	}
	return ret, nil
}

func NewRemoteCommand(caller Caller, cmdDef *CommandDefinition) *RemoteCommand {
	return &RemoteCommand{
		cmdDef: cmdDef,
		caller: caller,
	}
}

type CommandType int

const (
	NoKeyType CommandType = 1 << iota
	SingleKeyType
	MultiKeyType
)

type CommandDefinition struct {
	name string
	args []interface{}
	t    CommandType
}

func (self *CommandDefinition) Name() string {
	return self.name
}

func (self *CommandDefinition) Arg(i int) interface{} {
	return self.args[i]
}

func (self *CommandDefinition) Payload() Payload {
	payload := make([]interface{}, 1+len(self.args))
	payload[0] = self.name
	for i, a := range self.args {
		payload[i+1] = a
	}
	return payload
}

func (self *CommandDefinition) WithType(t CommandType) *CommandDefinition {
	self.t = t
	return self
}
func (self *CommandDefinition) Type() CommandType {
	return self.t
}

func (self *CommandDefinition) IsType(t CommandType) bool {
	return self.t&t != 0
}

func NewCommandDefinition(name string, args ...interface{}) *CommandDefinition {
	return &CommandDefinition{
		name: name,
		args: args,
		t:    SingleKeyType,
	}
}
