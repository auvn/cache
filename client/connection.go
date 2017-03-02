package client

import (
	"net"
	"time"

	"github.com/auvn/go.cache/net/serializer"
)

type Connection interface {
	Send(payload interface{}) error
	Receive() (serializer.Payload, error)
	Close() error
	Active() bool
	SetAuthenticated(auth bool)
	Authenticated() bool
}

type ConnectionOptions struct {
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

type connection struct {
	conn          net.Conn
	rw            serializer.ReadWriter
	inactive      bool
	authenticated bool
	opts          *ConnectionOptions
}

func (self *connection) checkActive(err error) {
	if nerr, ok := err.(net.Error); ok && !nerr.Temporary() {
		self.inactive = true
	}
}

func (self *connection) Send(payload interface{}) error {
	if err := self.rw.Write(payload); err != nil {
		self.checkActive(err)
		return err
	}
	return nil
}

func (self *connection) Receive() (serializer.Payload, error) {
	if p, err := self.rw.Read(); err != nil {
		self.checkActive(err)
		return nil, err
	} else {
		return p, nil
	}
}

func (self *connection) Close() error {
	return self.conn.Close()
}

func (self *connection) Active() bool {
	return !self.inactive
}

func (self *connection) SetAuthenticated(auth bool) {
	self.authenticated = auth
}

func (self *connection) Authenticated() bool {
	return self.authenticated
}

func NewConnection(conn net.Conn) Connection {
	return &connection{
		conn: conn,
		rw:   serializer.NewReadWriter(conn, conn),
	}
}

type ConnFactory interface {
	New() (Connection, error)
}

type ConnFactoryFunc func() (Connection, error)

func (self ConnFactoryFunc) New() (Connection, error) {
	return self()
}

type connFactory struct {
	addr    string
	timeout time.Duration
}

func (self *connFactory) New() (Connection, error) {
	conn, err := net.DialTimeout("tcp", self.addr, self.timeout)
	if err != nil {
		return nil, err
	}
	return NewConnection(conn), nil
}

func newConnectionFactory(addr string, timeout time.Duration) *connFactory {
	return &connFactory{addr: addr, timeout: timeout}
}

type Pool interface {
	Get() (*PooledConnection, error)
	Put(conn *PooledConnection)
}

type PooledConnection struct {
	Connection
	p Pool
}

func NewPooledConnection(p Pool, conn Connection) *PooledConnection {
	return &PooledConnection{
		p:          p,
		Connection: conn,
	}
}

type pool struct {
	conns       chan *PooledConnection
	connFactory ConnFactory
}

func (self *pool) Get() (*PooledConnection, error) {
	var c *PooledConnection
	select {
	case c = <-self.conns:
	default:
		rawConn, err := self.connFactory.New()
		if err != nil {
			return nil, err
		}
		c = NewPooledConnection(self, rawConn)
	}
	return c, nil
}

func (self *pool) Put(conn *PooledConnection) {
	if !conn.Active() {
		conn.Close()
	}
	select {
	case self.conns <- conn:
	default:
		conn.Close()
	}
}

func NewPool(capacity int, connFactory ConnFactory) Pool {
	if capacity < 0 {
		capacity = 0
	}
	return &pool{
		conns:       make(chan *PooledConnection, capacity),
		connFactory: connFactory,
	}
}
