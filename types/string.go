package types

import (
	"github.com/auvn/go.cache/core"
)

type String interface {
	Set(s core.Value)
	Get() core.Value
}

type strObject struct {
	value core.Value
}

func (self *strObject) Set(v core.Value) {
	self.value = v
}

func (self *strObject) Get() core.Value {
	return self.value
}

func NewString(v core.Value) String {
	return &strObject{value: v}
}
