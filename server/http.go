package server

import (
	"errors"
	"io"
	"log"
	"net/http"

	"strings"

	net "net"

	"github.com/auvn/go.cache/net/serializer"
	"github.com/auvn/go.cache/session"
	"github.com/auvn/go.cache/util/sync"
)

var (
	errEmptyPayload   = errors.New("empty payload")
	emptyPayload      = new(payload)
	emptyArrayPayload = make([]serializer.Payload, 0)
)

type payload struct{}

func (self *payload) Array() ([]serializer.Payload, error) {
	return emptyArrayPayload, nil
}

func (self *payload) Bytes() ([]byte, error) {
	return make([]byte, 0), errEmptyPayload
}

func (self *payload) Str() (string, error) {
	return "", errEmptyPayload
}

func (self *payload) Int() (int, error) {
	return 0, errEmptyPayload
}

func (self *payload) Bool() (bool, error) {
	return false, errEmptyPayload
}

func (self *payload) Err() error {
	return errEmptyPayload
}

func (self *payload) IsNil() bool {
	return true
}

func (self *payload) IsArray() bool {
	return true
}

func (self *payload) IsErr() bool {
	return false
}

type HTTPHandler struct {
	handler Handler
	session session.Session
	quit    sync.Quit
}

func (self *HTTPHandler) normalizePath(path string) string {
	if len(path) <= 0 {
		return path
	}
	if path[0] != '/' {
		return path
	}
	return strings.ToUpper(path[1:])
}

func (self *HTTPHandler) authHTTPRequest(s session.Session, req *http.Request, w serializer.Writer) error {
	if s.Authenticated() {
		return nil
	}
	password, _, _ := req.BasicAuth()
	authReq, err := self.prepareAuthRequest(password, s)
	if err != nil {
		return err
	}

	_, err = self.serveRequest(authReq)
	if err != nil {
		return err
	}
	return nil
}

func (self *HTTPHandler) readArray(r serializer.Reader) ([]serializer.Payload, error) {
	payload, err := r.ReadArray()
	if err != nil {
		if err == io.EOF {
			payload = emptyPayload
		} else {
			return emptyArrayPayload, ErrBadRequest.WithCause(err)
		}
	}

	arguments, err := payload.Array()
	if err != nil {
		return emptyArrayPayload, ErrBadRequest.WithCause(err)
	}
	return arguments, nil
}

func (self *HTTPHandler) serveHTTPRequest(s session.Session, req *http.Request, w serializer.Writer) error {
	err := self.authHTTPRequest(s, req, w)
	if err != nil {
		return ErrForbidden.WithCause(err)
	}

	path := self.normalizePath(req.URL.Path)
	r := serializer.NewReader(req.Body)
	arguments, err := self.readArray(r)
	if err != nil {
		return err
	}

	request, err := self.prepareRequest(path, arguments, s)
	if err != nil {
		return ErrBadRequest.WithCause(err)
	}

	value, err := self.serveRequest(request)
	if err != nil {
		return ErrBadRequest.WithCause(err)
	}

	if err = w.Write(value); err != nil {
		return ErrInternal.WithCause(err)
	}
	return nil
}

func (self *HTTPHandler) prepareRequest(cmdName string, arguments []serializer.Payload, s session.Session) (*Request, error) {
	body := make([][]byte, 0, len(arguments)+1)
	if cmdName != "" {
		body = append(body, []byte(cmdName))
	}
	for _, a := range arguments {
		bs, err := a.Bytes()
		if err != nil {
			return nil, err
		}
		body = append(body, bs)
	}
	return NewRequest(body, s), nil
}

func (self *HTTPHandler) prepareAuthRequest(password string, s session.Session) (*Request, error) {
	return NewRequest([][]byte{[]byte("AUTH"), []byte(password)}, s), nil
}

func (self *HTTPHandler) serveRequest(req *Request) (interface{}, error) {
	return handleRequest(self.handler, req, self.quit)
}

func (self *HTTPHandler) writeError(w serializer.Writer, err error) {
	if e := w.Write(err); e != nil {
		log.Println("cannot write error:", e)
	}
}
func (self *HTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream")

	writer := serializer.NewWriter(w)
	s := session.WithAuth(self.session)
	err := self.serveHTTPRequest(s, req, writer)
	if err != nil {
		if serr, ok := err.(Error); ok && !serr.Internal() {
			w.WriteHeader(serr.Code())
			self.writeError(writer, serr.Cause())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}
}

type HttpOptions struct {
	Addr string
}

type HttpServer struct {
	*http.Server
	handlerQuit chan struct{}
}

func (self *HttpServer) listen() (net.Listener, error) {
	return net.Listen("tcp", self.Addr)

}

func (self *HttpServer) waitForQuit(quit sync.Quit, listener net.Listener) {
	go func() {
		select {
		case <-quit:
			close(self.handlerQuit)
			listener.Close()
		}
	}()
}

func (self *HttpServer) Serve(quit sync.Quit) error {
	listener, err := self.listen()
	if err != nil {
		return err
	}
	self.waitForQuit(quit, listener)
	self.Server.Serve(listener)
	return nil
}

func Http(handler Handler, session session.Session, opts *HttpOptions) *HttpServer {
	if opts == nil {
		opts = &HttpOptions{
			Addr: "localhost:1235",
		}
	}

	handlerQuit := make(chan struct{}, 1)
	return &HttpServer{
		handlerQuit: handlerQuit,
		Server: &http.Server{
			Addr:    opts.Addr,
			Handler: NewHTTPHandler(handler, session, handlerQuit),
		},
	}
}

func NewHTTPHandler(handler Handler, session session.Session, quit sync.Quit) *HTTPHandler {
	return &HTTPHandler{
		session: session,
		handler: handler,
		quit:    quit,
	}
}
