package storage

import (
	"time"

	"github.com/auvn/go.cache/core"
)

const (
	MaxTTL time.Duration = 1<<63 - 1
)

type RawStorage interface {
	Get(ket core.StrValue) interface{}
	Set(key core.StrValue, v interface{})
	Del(key core.StrValue) bool
	TTL(key core.StrValue) core.IntValue
	SetTTL(key core.StrValue, ttl core.IntValue) bool
	Keys() []core.StrValue
	TimeNow() time.Time
}

type rawStorage struct {
	m map[core.StrValue]*ValueObject
	h *TTLHeap
}

func (self *rawStorage) del(key core.StrValue) {
	delete(self.m, key)
}

func (self *rawStorage) get(key core.StrValue, checkExpired bool) *ValueObject {
	v, ok := self.m[key]
	if !ok {
		return nil
	}

	if checkExpired && v.Expired(self.TimeNow()) {
		return nil
	}

	return v
}

func (self *rawStorage) Get(key core.StrValue) interface{} {
	if v := self.get(key, true); v != nil {
		return v.Object
	} else {
		return nil
	}
}

func (self *rawStorage) TTL(key core.StrValue) core.IntValue {
	var ret core.IntValue
	if v := self.get(key, true); v != nil {
		if v.Deadline().IsZero() {
			ret = -1
		} else {
			ret = core.IntValue(v.Deadline().Sub(self.TimeNow()).Seconds())
		}
	} else {
		ret = -2
	}
	return ret
}

func (self *rawStorage) Set(key core.StrValue, v interface{}) {
	self.Cleanup()
	self.m[key] = NewValueObject(v)
}

func (self *rawStorage) SetTTL(key core.StrValue, ttl core.IntValue) bool {
	defer self.Cleanup()

	v := self.get(key, true)
	if v == nil {
		return false
	}

	ttlValue := ttl.Value()

	if ttlValue < 0 {
		ttlValue = 0
	}

	ttlDuration := time.Duration(ttlValue) * time.Second
	if ttlDuration < 0 {
		ttlDuration = MaxTTL
	}
	deadline := self.TimeNow().Add(ttlDuration)
	if v.UpdateDeadline(deadline) {
		self.h.Fix(v)
	} else {
		self.h.Push(key, v)
	}
	return true
}

func (self *rawStorage) Del(key core.StrValue) bool {
	if v := self.get(key, true); v != nil {
		self.del(key)
		self.h.Delete(v)
		return true
	}
	return false
}

func (self *rawStorage) TimeNow() time.Time {
	return time.Now().UTC()
}

func (self *rawStorage) Keys() []core.StrValue {
	keys := make([]core.StrValue, 0, len(self.m))
	for k, v := range self.m {
		if !v.Expired(self.TimeNow()) {
			keys = append(keys, k)
		}
	}
	return keys
}

func (self *rawStorage) Cleanup() {
	for {
		if key, ok := self.h.PopExpired(self.TimeNow()); ok {
			// making sure the heap has fresh information about the key
			if v := self.get(key, false); v != nil && v.Expired(self.TimeNow()) {
				self.Del(key)
			}
		} else {
			break
		}
	}
}

type WriteFn func(Writer) (interface{}, error)
type ReadFn func(Reader) (interface{}, error)

type Storage interface {
	Write(fn WriteFn) (interface{}, error)
	Read(fn ReadFn) (interface{}, error)
}

type BaseStorage struct {
	reader Reader
	writer Writer
}

func (self *BaseStorage) Write(fn WriteFn) (interface{}, error) {
	return fn(self.writer)
}

func (self *BaseStorage) Read(fn ReadFn) (interface{}, error) {
	return fn(self.reader)
}

func New() Storage {
	rawStorage := &rawStorage{
		m: map[core.StrValue]*ValueObject{},
		h: NewTTLHeap(),
	}
	reader := &reader{storage: rawStorage}
	writer := &writer{Reader: reader, storage: rawStorage}
	return &BaseStorage{
		reader: reader,
		writer: writer,
	}
}
