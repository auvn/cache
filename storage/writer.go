package storage

import (
	"github.com/auvn/go.cache/core"
)

type Writer interface {
	Reader
	Set(key core.StrValue, v interface{})
	SetTTL(key core.StrValue, ttl core.IntValue) bool
	Delete(key core.StrValue) bool
}

type writer struct {
	Reader
	storage RawStorage
}

func (self *writer) Set(key core.StrValue, v interface{}) {
	self.storage.Set(key, v)
}

func (self *writer) SetTTL(key core.StrValue, ttl core.IntValue) bool {
	return self.storage.SetTTL(key, ttl)
}

func (self *writer) Delete(key core.StrValue) bool {
	return self.storage.Del(key)
}
