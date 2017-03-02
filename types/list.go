package types

import (
	"github.com/auvn/go.cache/core"
)

var (
	emptyValue      = core.Value{}
	emptyValueSlice = []core.Value{}
)

type List interface {
	LPush(values ...core.Value) core.IntValue
	RPush(values ...core.Value) core.IntValue
	LPop() (core.Value, bool)
	RPop() (core.Value, bool)
	Range(start, stop core.IntValue) []core.Value
	Get(index core.IntValue) core.Value
}

type listElement struct {
	Value core.Value
	Next  *listElement
	Prev  *listElement
}

type listObject struct {
	Head   *listElement
	Tail   *listElement
	Length int
}

func (self *listObject) push(beginning bool, values ...core.Value) core.IntValue {
	if len(values) == 0 {
		return self.length()
	}

	if self.empty() {
		self.pushFirst(values[0])
		values = values[1:]
	}
	if beginning {
		for _, value := range values {
			self.lpush(value)
		}
	} else {
		for _, value := range values {
			self.rpush(value)
		}
	}
	return self.length()
}

func (self *listObject) LPush(values ...core.Value) core.IntValue {
	return self.push(true, values...)
}

func (self *listObject) RPush(values ...core.Value) core.IntValue {
	return self.push(false, values...)
}

func (self *listObject) pop(beginning bool) (core.Value, bool) {
	if self.empty() {
		return emptyValue, false
	}

	if beginning {
		return self.lpop().Value, true
	} else {
		return self.rpop().Value, true
	}
}

func (self *listObject) LPop() (core.Value, bool) {
	return self.pop(true)
}

func (self *listObject) RPop() (core.Value, bool) {
	return self.pop(false)
}

func (self *listObject) Range(start, stop core.IntValue) []core.Value {
	startVal := start.Value()
	stopVal := stop.Value()
	length := self.Length

	if startVal < 0 {
		startVal = length + startVal
	}
	if startVal < 0 {
		startVal = 0
	}
	if stopVal < 0 {
		stopVal = length + stopVal
	}

	if startVal > stopVal || startVal >= length {
		return emptyValueSlice
	}

	if stopVal >= length {
		stopVal = length - 1
	}
	cursor := self.Head

	for i := 0; i < startVal; i++ {
		cursor = cursor.Next
	}
	count := stopVal - startVal
	values := make([]core.Value, 0, count)
	for count >= 0 {
		values = append(values, cursor.Value)
		count -= 1
		cursor = cursor.Next
	}

	return values
}

func (self *listObject) Get(index core.IntValue) core.Value {
	i := index.Value()
	if self.Length < i {
		return emptyValue
	}

	var cur int
	elem := self.Head
	for cur < i && elem != nil {
		elem = elem.Next
		cur += 1
	}
	if elem == nil {
		return emptyValue
	}
	return elem.Value
}

func (self *listObject) length() core.IntValue {
	return core.IntValue(self.Length)
}

func (self *listObject) empty() bool {
	return self.Length == 0
}

func (self *listObject) adjustLength(value int) int {
	self.Length += value
	return self.Length
}

func (self *listObject) pushFirst(value core.Value) *listElement {
	defer self.adjustLength(1)

	elem := &listElement{Value: value, Next: nil, Prev: nil}
	self.Head = elem
	self.Tail = elem
	return elem
}

func (self *listObject) lpush(value core.Value) *listElement {
	defer self.adjustLength(1)

	oldHead := self.Head
	self.Head = &listElement{Value: value, Next: oldHead, Prev: nil}
	oldHead.Prev = self.Head
	return self.Head
}

func (self *listObject) rpush(value core.Value) *listElement {
	defer self.adjustLength(1)

	oldTail := self.Tail
	self.Tail = &listElement{Value: value, Prev: oldTail, Next: nil}
	oldTail.Next = self.Tail
	return self.Tail
}

func (self *listObject) lpop() *listElement {
	defer self.adjustLength(-1)

	elem := self.Head
	self.Head = elem.Next
	if self.Head != nil {
		self.Head.Prev = nil
	}

	elem.Next = nil
	elem.Prev = nil
	return elem
}

func (self *listObject) rpop() *listElement {
	defer self.adjustLength(-1)

	elem := self.Tail
	self.Tail = elem.Prev
	if self.Tail != nil {
		self.Tail.Next = nil
	}

	elem.Next = nil
	elem.Prev = nil
	return elem
}

func NewList() List {
	return &listObject{}
}
