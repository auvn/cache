package server

import (
	"errors"

	"github.com/auvn/go.cache/session"
	"github.com/auvn/go.cache/util/sync"
)

var (
	ErrQuit = errors.New("quit")
)

type response chan interface{}

func (self response) Get(quit sync.Quit) (interface{}, error) {
	select {
	case v := <-self:
		if err, ok := v.(error); ok {
			return nil, err
		} else {
			return v, nil
		}
	case <-quit:
		return nil, ErrQuit
	}
}

func (self response) Set(v interface{}) {
	self <- v
}

type Request struct {
	body     [][]byte
	session  session.Session
	response response
}

func (self *Request) Body() [][]byte {
	return self.body
}

func (self *Request) Session() session.Session {
	return self.session
}

func (self *Request) Response() chan<- interface{} {
	return self.response
}

func NewRequest(body [][]byte, s session.Session) *Request {
	return &Request{
		body:     body,
		session:  s,
		response: make(response, 1),
	}
}

type Handler interface {
	HandleRequest(*Request)
}

func handleRequest(handler Handler, req *Request, quit sync.Quit) (interface{}, error) {
	handler.HandleRequest(req)
	return req.response.Get(quit)
}
