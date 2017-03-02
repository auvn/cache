package storage

import (
	"container/heap"
	"time"

	"github.com/auvn/go.cache/core"
)

type Indexable interface {
	Index() int
	SetIndex(int)
}

type Expirable interface {
	Indexable
	Expired(time.Time) bool
	Deadline() time.Time
}

type ttlQueue []Expirable

func (self ttlQueue) Len() int {
	return len(self)
}

func (self ttlQueue) Less(i, j int) bool {
	id := self[i].Deadline()
	jd := self[j].Deadline()
	return id.Before(jd)
}

func (self ttlQueue) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
	self[i].SetIndex(j)
	self[j].SetIndex(i)
}

func (self *ttlQueue) Push(x interface{}) {
	k := x.(Expirable)
	index := len(*self)
	*self = append(*self, k)
	k.SetIndex(index)
}

func (self *ttlQueue) Pop() interface{} {
	h := *self
	hlen := len(h)
	k := h[hlen-1]
	*self = h[:hlen-1]
	k.SetIndex(1)
	return k
}

type keyTTL struct {
	Expirable
	key core.StrValue
}

type TTLHeap struct {
	h *ttlQueue
}

func (self *TTLHeap) Push(key core.StrValue, v Expirable) {
	k := &keyTTL{key: key, Expirable: v}
	heap.Push(self.h, k)
}

func (self *TTLHeap) Pop() core.StrValue {
	v := heap.Pop(self.h)
	if v == nil {
		return core.EmptyStrValue
	} else {
		return v.(*keyTTL).key
	}
}

func (self *TTLHeap) PopExpired(now time.Time) (core.StrValue, bool) {
	arr := *self.h
	if arr.Len() <= 0 {
		return core.EmptyStrValue, false
	}
	k := arr[0].(*keyTTL)
	if k.Expired(now) {
		self.Pop()
		return k.key, true
	}
	return core.EmptyStrValue, false
}

func (self *TTLHeap) Fix(i Indexable) {
	index := i.Index()
	if index < 0 {
		return
	}
	heap.Fix(self.h, index)
}

func (self *TTLHeap) Delete(i Indexable) {
	index := i.Index()
	if index < 0 {
		return
	}
	heap.Remove(self.h, index)
}

func NewTTLHeap() *TTLHeap {
	h := &ttlQueue{}
	heap.Init(h)
	return &TTLHeap{h: h}
}
