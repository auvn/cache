package session

import (
	"errors"
	"sync"

	"github.com/auvn/go.cache/storage"
)

var (
	ErrEmptyStorage = errors.New("empty storage")

	Empty        = new(emptySession)
	emptyStorage = new(emptyStorageObj)
)

type Session interface {
	Authenticated() bool
	SetAuthenticated(bool)
	Storage() storage.Storage
}

type emptySession struct{}

func (self *emptySession) Authenticated() bool {
	return true
}

func (self *emptySession) SetAuthenticated(auth bool) {

}

func (self *emptySession) Storage() storage.Storage {
	return emptyStorage
}

func New() Session {
	return Empty
}

type authSession struct {
	Session
	rw   sync.RWMutex
	auth bool
}

func (self *authSession) Authenticated() bool {
	self.rw.RLock()
	defer self.rw.RUnlock()
	return self.auth
}

func (self *authSession) SetAuthenticated(auth bool) {
	self.rw.Lock()
	defer self.rw.Unlock()
	self.auth = auth
}

func WithAuth(s Session) Session {
	return &authSession{
		Session: s,
	}
}

type emptyStorageObj struct{}

func (self *emptyStorageObj) Write(fn storage.WriteFn) (interface{}, error) {
	return nil, ErrEmptyStorage
}

func (self *emptyStorageObj) Read(fn storage.ReadFn) (interface{}, error) {
	return nil, ErrEmptyStorage
}

type storageSession struct {
	Session
	s storage.Storage
}

func (self *storageSession) Storage() storage.Storage {
	return self.s
}

func WithStorage(s Session, stor storage.Storage) Session {
	return &storageSession{
		Session: s,
		s:       stor,
	}
}
