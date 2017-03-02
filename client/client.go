package client

import (
	"hash/crc32"
	"net"
	"sync"
	"time"

	"github.com/auvn/go.cache/net/serializer"
)

type Auther interface {
	Auth(Connection) error
}

type simpleCaller struct {
	conn Connection
}

func (self *simpleCaller) Call(cmdDef *CommandDefinition) (serializer.Payload, error) {
	return execute(self.conn, cmdDef.Payload())
}

func newSimpleCaller(conn Connection) *simpleCaller {
	return &simpleCaller{conn: conn}
}

type auther struct {
	cmdDef *CommandDefinition
}

func (self *auther) Auth(conn Connection) error {
	if conn.Authenticated() {
		return nil
	}
	caller := newSimpleCaller(conn)
	remoteCommand := NewRemoteCommand(caller, self.cmdDef)
	_, err := remoteCommand.Bool()
	if err != nil {
		return err
	}
	conn.SetAuthenticated(true)
	return nil
}

func newAuther(cmdDef *CommandDefinition) *auther {
	return &auther{cmdDef: cmdDef}
}

type dummyAuther struct{}

func (self *dummyAuther) Auth(conn Connection) error {
	return nil
}

type Client interface {
	Caller
}

func execute(conn Connection, payload Payload) (serializer.Payload, error) {
	err := conn.Send(payload)
	if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
		return nil, nerr
	} // else -- trying to read the last data from the connection
	// maybe:
	// store network non-temp error
	// make the last chance and read the conn,
	// on serializer.Error return the stored network error
	response, err := conn.Receive()
	if err != nil {
		return nil, err
	}

	if response.IsErr() {
		return nil, response.Err()
	}
	return response, nil
}

type baseClient struct {
	auther  Auther
	options *Options
	pool    Pool
}

func (self *baseClient) Call(cmdDef *CommandDefinition) (serializer.Payload, error) {
	conn, err := self.pool.Get()
	if err != nil {
		return nil, err
	}
	defer self.pool.Put(conn)

	if err = self.auther.Auth(conn); err != nil {
		return nil, err
	}
	return execute(conn, cmdDef.Payload())
}

func newBaseClient(addr string, poolSize int, dialTimeout time.Duration, auther Auther) *baseClient {
	return &baseClient{
		auther: auther,
		pool:   NewPool(poolSize, newConnectionFactory(addr, dialTimeout)),
	}
}

type multiClient struct {
	auther       Auther
	pools        []Pool
	serversCount int
	hash         func(key string) uint32
}

func (self *multiClient) poolIndex(key string) int {
	return int(self.hash(key) % uint32(self.serversCount))
}

func (self *multiClient) callAsync(index int, payload Payload, ch chan interface{}) {
	go func() {
		p, err := self.call(index, payload)
		if err != nil {
			ch <- err
		} else {
			ch <- p
		}
	}()
}

func (self *multiClient) multiCall(payload Payload) (serializer.Payload, error) {
	n := self.serversCount
	wg := &sync.WaitGroup{}
	wg.Add(n)
	responses := make(chan interface{}, n)
	for i, _ := range self.pools {
		self.callAsync(i, payload, responses)
	}

	go func(wg *sync.WaitGroup, ch chan interface{}) {
		wg.Wait()
		close(ch)
	}(wg, responses)

	results := make([]serializer.Payload, 0, n)
	for resp := range responses {
		wg.Done()
		if err, ok := resp.(error); ok {
			return nil, err
		} else {
			results = append(results, resp.(serializer.Payload))
		}
	}
	return MultiPayload(results), nil

}
func (self *multiClient) call(index int, payload Payload) (serializer.Payload, error) {
	pool := self.pools[index]
	conn, err := pool.Get()
	if err != nil {
		return nil, err
	}
	defer pool.Put(conn)
	if err = self.auther.Auth(conn); err != nil {
		return nil, err
	}
	return execute(conn, payload)
}

func (self *multiClient) Call(cmdDef *CommandDefinition) (serializer.Payload, error) {
	payload := cmdDef.Payload()
	if cmdDef.IsType(NoKeyType | MultiKeyType) {
		return self.multiCall(payload)
	} else {
		key := cmdDef.Arg(0).(string)
		i := self.poolIndex(key)
		return self.call(i, payload)
	}
}

func newMultiClient(addrs []string, poolSize int, dialTimeout time.Duration, auther Auther) *multiClient {
	serversCount := len(addrs)
	pools := make([]Pool, serversCount)
	for i, _ := range pools {
		pools[i] = NewPool(poolSize, newConnectionFactory(addrs[i], dialTimeout))
	}
	return &multiClient{
		auther:       auther,
		pools:        pools,
		serversCount: serversCount,
		hash: func(key string) uint32 {
			return crc32.ChecksumIEEE([]byte(key))
		},
	}
}
