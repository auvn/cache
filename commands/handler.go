package commands

import (
	"errors"

	"github.com/auvn/go.cache/core"
	"github.com/auvn/go.cache/server"
	"github.com/auvn/go.cache/session"
	"github.com/auvn/go.cache/util/sync"
)

var (
	ErrCannotUpdateJournal = errors.New("cannot perform an update in the journal")
)

type SuccessHook func(flag int, body [][]byte)

type Handler struct {
	registry     Registry
	requests     chan *server.Request
	successHooks []SuccessHook
}

func (self *Handler) lookupCommand(values []core.Value) (Command, Arguments, error) {
	arguments := NewArguments(values...)
	iter, err := arguments.IterAtLeast(1)
	if err != nil {
		return nil, nil, err
	}
	cmdName, err := iter.NextStr()
	if err != nil {
		return nil, nil, err
	}

	cmd, ok := self.registry.Get(cmdName.Value())
	if !ok {
		return nil, nil, ErrUnknownCommand
	}
	return cmd, iter.NextArguments(), nil
}

func (self *Handler) handle(body [][]byte, s session.Session, resp chan<- interface{}) (Command, error) {
	values := make([]core.Value, len(body))
	for i, _ := range body {
		values[i] = body[i]
	}
	cmd, arguments, err := self.lookupCommand(values)
	if err != nil {
		resp <- err
		return nil, err
	}
	ret, err := cmd.Execute(s, arguments)
	if err != nil {
		resp <- err
	} else {
		resp <- ret
	}
	return cmd, err
}

func (self *Handler) handleRequest(req *server.Request) {
	body := req.Body()
	sess := req.Session()
	resp := req.Response()
	cmd, err := self.handle(body, sess, resp)
	if err == nil {
		for _, h := range self.successHooks {
			h(cmd.Flag(), body)
		}
	}
}

func (self *Handler) loopRequests(quit sync.Quit) {
	for {
		select {
		case <-quit:
			return
		case req := <-self.requests:
			self.handleRequest(req)
		}
	}
}

func (self *Handler) Serve(quit sync.Quit) error {
	self.loopRequests(quit)
	return nil
}

func (self *Handler) Handle(s session.Session, body [][]byte, resp chan<- interface{}) {
	self.handle(body, s, resp)
}

func (self *Handler) HandleRequest(req *server.Request) {
	self.requests <- req
}

func (self *Handler) AddSuccessHook(fn SuccessHook) {
	self.successHooks = append(self.successHooks, fn)
}

func NewHandler(registry Registry) *Handler {
	return &Handler{
		registry:     registry,
		requests:     make(chan *server.Request, 100),
		successHooks: make([]SuccessHook, 0, 10),
	}
}
