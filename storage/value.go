package storage

import (
	"time"
)

type ValueObject struct {
	Object   interface{}
	deadline time.Time
	index    int
}

func (self *ValueObject) Expired(t time.Time) bool {
	deadline := self.deadline
	return !deadline.IsZero() && deadline.Before(t)
}

// false if deadline was set to Zero
func (self *ValueObject) UpdateDeadline(deadline time.Time) bool {
	fresh := self.deadline.IsZero()
	self.deadline = deadline
	return !fresh
}

func (self *ValueObject) Deadline() time.Time {
	return self.deadline
}

func (self *ValueObject) Index() int {
	return self.index
}
func (self *ValueObject) SetIndex(i int) {
	self.index = i
}

func (self *ValueObject) SetObject(object interface{}) {
	self.Object = object
}

func NewValueObject(object interface{}) *ValueObject {
	return &ValueObject{Object: object, index: -1}
}
