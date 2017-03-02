package storage

import (
	"github.com/auvn/go.cache/core"
)

type Reader interface {
	Get(key core.StrValue) (interface{}, bool)
	TTL(key core.StrValue) core.IntValue
	Keys() []core.StrValue
}

type reader struct {
	storage RawStorage
}

func (self *reader) Get(key core.StrValue) (interface{}, bool) {
	if v := self.storage.Get(key); v != nil {
		return v, true
	}
	return nil, false
}

func (self *reader) TTL(key core.StrValue) core.IntValue {
	return self.storage.TTL(key)
}

func (self *reader) Keys() []core.StrValue {
	return self.storage.Keys()
}
