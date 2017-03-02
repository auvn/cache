package server

import (
	"net"
	"time"

	"github.com/auvn/go.cache/net/serializer"
	"github.com/auvn/go.cache/session"
	"github.com/auvn/go.cache/util/sync"
)

var (
	DefaultOptions = &TelnetOptions{
		Addr: ":1234",
	}
	ZeroTime = time.Time{}
)

type TelnetOptions struct {
	Addr string
}

type TelnetServer struct {
	opts    *TelnetOptions
	session session.Session
	group   sync.ServeGroup
	handler Handler
}

func (self *TelnetServer) listen() (net.Listener, error) {
	return net.Listen("tcp", self.opts.Addr)
}

func (self *TelnetServer) serveClient(conn net.Conn) {
	client := NewTelnetClient(
		self.handler,
		conn,
		session.WithAuth(self.session),
	)
	self.group.Serve(client)
}

func (self *TelnetServer) waitForQuit(listener net.Listener, quit sync.Quit) {
	go func() {
		select {
		case <-quit:
			listener.Close()
		}
	}()
}

func (self *TelnetServer) loopListener(listener net.Listener, quit sync.Quit) {
	for {
		select {
		case <-quit:
			self.group.Shutdown()
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
					time.Sleep(100 * time.Millisecond)
					continue
				} else {
					time.Sleep(1 * time.Second)
					continue
				}
			}
			self.serveClient(conn)
		}
	}
}

func (self *TelnetServer) Serve(quit sync.Quit) error {
	listener, err := self.listen()
	if err != nil {
		return err
	}
	self.waitForQuit(listener, quit)
	self.loopListener(listener, quit)
	return nil
}

func Telnet(handler Handler, session session.Session, opts *TelnetOptions) *TelnetServer {
	if opts == nil {
		opts = DefaultOptions
	}
	return &TelnetServer{
		handler: handler,
		session: session,
		opts:    opts,
		group:   sync.NewServeGroup(),
	}
}

var (
	DefaultTelnetClientOptions = &TelnetClientOptions{
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Minute,
	}
)

type TelnetClientOptions struct {
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

type TelnetClient struct {
	rw      serializer.ReadWriter
	conn    net.Conn
	session session.Session
	handler Handler
	opts    *TelnetClientOptions
}

func (self *TelnetClient) timeNow() time.Time {
	return time.Now()
}

func (self *TelnetClient) nextWriteDeadline() time.Time {
	return self.timeNow().Add(self.opts.WriteTimeout)
}

func (self *TelnetClient) nextReadDeadline() time.Time {
	return self.timeNow().Add(self.opts.ReadTimeout)
}

func (self *TelnetClient) write(i interface{}) {
	self.conn.SetWriteDeadline(self.nextWriteDeadline())
	defer self.conn.SetWriteDeadline(ZeroTime)
	self.rw.Write(i)
}

func (self *TelnetClient) readArray() (serializer.Payload, error) {
	self.conn.SetReadDeadline(self.nextReadDeadline())
	defer self.conn.SetReadDeadline(ZeroTime)
	return self.rw.ReadArray()
}

func (self *TelnetClient) waitForQuit(quit sync.Quit) {
	go func() {
		select {
		case <-quit:
			self.conn.Close()
		}
	}()
}

func (self *TelnetClient) loopCommands(quit sync.Quit) {
	defer self.conn.Close()
	for {
		select {
		case <-quit:
			return
		default:
			payload, err := self.readArray()
			if err != nil {
				if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
					continue
				} else if serr, ok := err.(serializer.Error); ok {
					self.write(serr)
					return
				} else {
					return
				}
			}

			req, err := self.prepareRequest(payload, self.session)
			if err != nil {
				self.write(err)
				continue
			}

			value, err := handleRequest(self.handler, req, quit)
			if err != nil {
				self.write(err)
			} else {
				self.write(value)
			}
		}
	}
}

func (self *TelnetClient) Serve(quit sync.Quit) error {
	self.waitForQuit(quit)
	self.loopCommands(quit)
	return nil
}

func (self *TelnetClient) prepareRequest(p serializer.Payload, s session.Session) (*Request, error) {
	array, err := p.Array()
	if err != nil {
		return nil, err
	}

	body := make([][]byte, len(array))
	for i, p := range array {
		value, err := p.Bytes()
		if err != nil {
			return nil, err
		}
		body[i] = value
	}

	return NewRequest(body, s), nil
}

func NewTelnetClient(handler Handler, conn net.Conn, session session.Session) *TelnetClient {
	return &TelnetClient{
		rw:      serializer.NewReadWriter(conn, conn),
		handler: handler,
		conn:    conn,
		session: session,
		opts:    DefaultTelnetClientOptions,
	}
}
