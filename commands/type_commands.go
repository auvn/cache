package commands

import (
	"errors"

	"github.com/auvn/go.cache/core"
	"github.com/auvn/go.cache/session"
	"github.com/auvn/go.cache/storage"
	"github.com/auvn/go.cache/types"
)

var (
	ErrForbidden = errors.New("forbidden")
)

type SecurityCommand string

func (self *SecurityCommand) Auth(s session.Session, pass core.StrValue) (interface{}, error) {
	password := string(*self)
	if s.Authenticated() {
		return true, nil
	} else if password == "" || pass.Value() == password {
		s.SetAuthenticated(true)
		return true, nil
	}
	return false, ErrForbidden
}

func NewSecurityCommand(pass string) *SecurityCommand {
	cmd := SecurityCommand(pass)
	return &cmd
}

type StorageCommand struct{}

func (self *StorageCommand) Del(s session.Session, keys ...core.StrValue) (interface{}, error) {
	return s.Storage().Write(func(w storage.Writer) (interface{}, error) {
		var counter int
		for _, k := range keys {
			if w.Delete(k) {
				counter += 1
			}
		}
		return counter, nil
	})
}

func (self *StorageCommand) Expire(s session.Session, key core.StrValue, ttl core.IntValue) (interface{}, error) {
	return s.Storage().Write(func(w storage.Writer) (interface{}, error) {
		return w.SetTTL(key, ttl), nil
	})
}
func (self *StorageCommand) TTL(s session.Session, key core.StrValue) (interface{}, error) {
	return s.Storage().Read(func(r storage.Reader) (interface{}, error) {
		return r.TTL(key), nil
	})
}

func (self *StorageCommand) Keys(s session.Session) (interface{}, error) {
	return s.Storage().Read(
		func(r storage.Reader) (interface{}, error) {
			return r.Keys(), nil
		},
	)
}

func NewStorageCommand() *StorageCommand {
	return new(StorageCommand)
}

type StringCommand struct{}

func (self *StringCommand) cast(v interface{}) (types.String, error) {
	if s, ok := v.(types.String); ok {
		return s, nil
	} else {
		return nil, ErrWrongType
	}
}

func (self *StringCommand) Set(s session.Session, key core.StrValue, value core.Value) (interface{}, error) {
	return s.Storage().Write(
		func(w storage.Writer) (interface{}, error) {
			w.Set(key, types.NewString(value))
			return true, nil
		},
	)
}

func (self *StringCommand) Get(s session.Session, key core.StrValue) (interface{}, error) {
	return s.Storage().Read(
		func(r storage.Reader) (interface{}, error) {
			if value, ok := r.Get(key); ok {
				if str, err := self.cast(value); err == nil {
					return str.Get(), nil
				} else {
					return nil, err
				}
			}
			return nil, nil
		},
	)
}

func NewStringCommand() *StringCommand {
	return new(StringCommand)
}

type HashCommand struct{}

func (self *HashCommand) cast(v interface{}) (types.Hash, error) {
	if h, ok := v.(types.Hash); ok {
		return h, nil
	} else {
		return nil, ErrWrongType
	}
}

func (self *HashCommand) Set(s session.Session, key core.StrValue, hashKey core.StrValue, hashValue core.Value) (interface{}, error) {
	return s.Storage().Write(func(w storage.Writer) (interface{}, error) {
		var h types.Hash
		var err error
		value, ok := w.Get(key)
		if ok {
			if h, err = self.cast(value); err != nil {
				return nil, err
			}
		} else {
			h = types.NewHash()
			w.Set(key, h)
		}
		return h.Set(hashKey, hashValue), nil
	})

}

func (self *HashCommand) Get(s session.Session, key core.StrValue, hashKey core.StrValue) (interface{}, error) {
	return s.Storage().Read(func(r storage.Reader) (interface{}, error) {
		var h types.Hash
		var err error

		if value, ok := r.Get(key); ok {
			if h, err = self.cast(value); err == nil {
				if v, ok := h.Get(hashKey); ok {
					return v, nil
				}
			}
		}
		return nil, nil
	})
}

func (self *HashCommand) Del(s session.Session, key core.StrValue, hashKeys ...core.StrValue) (interface{}, error) {
	return s.Storage().Write(
		func(w storage.Writer) (interface{}, error) {
			if value, ok := w.Get(key); ok {
				if h, err := self.cast(value); err == nil {
					return h.Del(hashKeys...), nil
				}
			}
			return nil, nil
		},
	)
}

func (self *HashCommand) Keys(s session.Session, key core.StrValue) (interface{}, error) {
	return s.Storage().Read(
		func(r storage.Reader) (interface{}, error) {
			if value, ok := r.Get(key); ok {
				if h, err := self.cast(value); err == nil {
					return h.Keys(), nil
				}
			}
			return nil, nil
		},
	)
}

func NewHashCommand() *HashCommand {
	return new(HashCommand)
}

type ListCommand struct{}

func (self *ListCommand) cast(v interface{}) (types.List, error) {
	if s, ok := v.(types.List); ok {
		return s, nil
	} else {
		return nil, ErrWrongType
	}
}

func (self *ListCommand) push(s session.Session, beginning bool, key core.StrValue, values ...core.Value) (interface{}, error) {
	return s.Storage().Write(func(w storage.Writer) (interface{}, error) {
		var err error
		var l types.List
		value, ok := w.Get(key)
		if ok {
			if l, err = self.cast(value); err != nil {
				return nil, err
			}
		} else {
			l = types.NewList()
			w.Set(key, l)
		}
		if beginning {
			return l.LPush(values...), nil
		} else {
			return l.LPush(values...), nil
		}
	})
}

func (self *ListCommand) pop(s session.Session, beginning bool, key core.StrValue) (interface{}, error) {
	return s.Storage().Write(func(w storage.Writer) (interface{}, error) {
		var err error
		var l types.List
		if value, ok := w.Get(key); ok {
			if l, err = self.cast(value); err != nil {
				return nil, err
			} else {
				if beginning {
					v, _ := l.LPop()
					return v, nil
				} else {
					v, _ := l.RPop()
					return v, nil
				}
			}
		}
		return nil, nil
	})
}

func (self *ListCommand) LPush(s session.Session, key core.StrValue, values ...core.Value) (interface{}, error) {
	return self.push(s, true, key, values...)
}

func (self *ListCommand) RPush(s session.Session, key core.StrValue, values ...core.Value) (interface{}, error) {
	return self.push(s, false, key, values...)
}

func (self *ListCommand) LPop(s session.Session, key core.StrValue) (interface{}, error) {
	return self.pop(s, true, key)
}

func (self *ListCommand) RPop(s session.Session, key core.StrValue) (interface{}, error) {
	return self.pop(s, false, key)
}

func (self *ListCommand) LRange(s session.Session, key core.StrValue, start, stop core.IntValue) (interface{}, error) {
	return s.Storage().Read(func(r storage.Reader) (interface{}, error) {
		var l types.List
		var err error
		if value, ok := r.Get(key); ok {
			if l, err = self.cast(value); err != nil {
				return nil, err
			}
			return l.Range(start, stop), nil
		}
		return nil, nil
	})
}

func (self *ListCommand) LIndex(s session.Session, key core.StrValue, index core.IntValue) (interface{}, error) {
	return s.Storage().Read(func(r storage.Reader) (interface{}, error) {
		var err error
		var l types.List
		if value, ok := r.Get(key); ok {
			if l, err = self.cast(value); err != nil {
				return nil, err
			}
			return l.Get(index), nil
		}
		return nil, nil
	})
}

func NewListCommand() *ListCommand {
	return new(ListCommand)
}
