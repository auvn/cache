package types

import "github.com/auvn/go.cache/core"

type Hash interface {
	Set(key core.StrValue, value core.Value) bool
	Get(key core.StrValue) (core.Value, bool)
	Del(key ...core.StrValue) core.IntValue
	Keys() []core.StrValue
}

type hashStorage map[core.StrValue]core.Value

func (self hashStorage) Get(key core.StrValue) (core.Value, bool) {
	value, ok := self[key]
	return value, ok
}

func (self hashStorage) Set(key core.StrValue, value core.Value) {
	self[key] = value
}

func (self hashStorage) Delete(key core.StrValue) {
	delete(self, key)
}

type hashObject struct {
	storage hashStorage
}

func (self *hashObject) Set(key core.StrValue, value core.Value) bool {
	_, updated := self.storage.Get(key)
	self.storage.Set(key, value)
	return !updated
}

func (self *hashObject) Get(key core.StrValue) (core.Value, bool) {
	return self.storage.Get(key)
}

func (self *hashObject) Del(keys ...core.StrValue) core.IntValue {
	var counter int = 0
	for _, k := range keys {
		_, ok := self.storage.Get(k)
		if !ok {
			continue
		}
		self.storage.Delete(k)
		counter += 1
	}
	return core.IntValue(counter)
}

func (self *hashObject) Keys() []core.StrValue {
	keys := make([]core.StrValue, 0, len(self.storage))
	for k := range self.storage {
		keys = append(keys, k)
	}
	return keys
}

func NewHash() Hash {
	return &hashObject{storage: hashStorage{}}
}
