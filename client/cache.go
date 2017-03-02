package client

import "time"

const (
	AuthCommand   = "AUTH"
	DelCommand    = "DEL"
	KeysCommand   = "KEYS"
	TTLCommand    = "TTL"
	ExpireCommand = "EXPIRE"

	//string
	SetCommand = "SET"
	GetCommand = "GET"

	//list
	LPushCommand  = "LPUSH"
	RPushCommand  = "RPUSH"
	LPopCommand   = "LPOP"
	RPopCommand   = "RPOP"
	LRangeCommand = "LRANGE"
	LIndexCommand = "LINDEX"

	//hash
	HSetCommand  = "HSET"
	HGetCommand  = "HGET"
	HKeysCommand = "HKEYS"
	HDelCommand  = "HDEL"
)

var (
	DefaultOptions = &Options{
		Addrs:       []string{":1234"},
		Auth:        "",
		PoolSize:    10,
		DialTimeout: 5 * time.Second,
	}
)

type Options struct {
	Addrs       []string
	Auth        string
	PoolSize    int
	DialTimeout time.Duration
}

type Cache interface {
	Del(keys ...string) IntCommand
	Keys() StringSliceCommand
	TTL(key string) IntCommand
	Expire(key string, ttl int) BoolCommand

	Get(key string) BytesCommand
	Set(key string, value []byte) BoolCommand

	LIndex(key string) BytesCommand
	LPop(key string) BytesCommand
	LPush(key string, values ...[]byte) IntCommand
	LRange(key string, start int, stop int) BytesSliceCommand
	RPop(key string) BytesCommand
	RPush(key string, values ...[]byte) IntCommand

	HDel(key string, hashKeys ...[]byte) IntCommand
	HGet(key string, hashKey []byte) BytesCommand
	HKeys(key string) StringSliceCommand
	HSet(key string, hashKey []byte, value []byte) BoolCommand
}

type cache struct {
	client Client
}

func (self *cache) command(cmdDef *CommandDefinition) Command {
	return NewRemoteCommand(self.client, cmdDef)
}

func (self *cache) Del(keys ...string) IntCommand {
	args := make([]interface{}, len(keys))
	for i, k := range keys {
		args[i] = k
	}
	cmdDef := NewCommandDefinition(DelCommand, args...).WithType(MultiKeyType)
	return self.command(cmdDef)
}

func (self *cache) Keys() StringSliceCommand {
	return self.command(
		NewCommandDefinition(KeysCommand).WithType(NoKeyType),
	)
}

func (self *cache) TTL(key string) IntCommand {
	cmdDef := NewCommandDefinition(TTLCommand, key)
	return self.command(cmdDef)
}

func (self *cache) Expire(key string, ttl int) BoolCommand {
	cmdDef := NewCommandDefinition(ExpireCommand, key, ttl)
	return self.command(cmdDef)
}

///////////////////////// string ////////////////////////
func (self *cache) Set(key string, value []byte) BoolCommand {
	cmdDef := NewCommandDefinition(SetCommand, key, value)
	return self.command(cmdDef)
}

func (self *cache) Get(key string) BytesCommand {
	cmdDef := NewCommandDefinition(GetCommand, key)
	return self.command(cmdDef)
}

///////////////////////// list ////////////////////////
func (self *cache) Push(beginning bool, key string, values ...[]byte) IntCommand {
	var cmdName string
	if beginning {
		cmdName = LPushCommand
	} else {
		cmdName = RPushCommand
	}
	args := make([]interface{}, 1+len(values))
	args[0] = key
	for i, v := range values {
		args[i+1] = v
	}
	cmdDef := NewCommandDefinition(cmdName, args...)
	return self.command(cmdDef)
}

func (self *cache) LPush(key string, values ...[]byte) IntCommand {
	return self.Push(true, key, values...)
}

func (self *cache) RPush(key string, values ...[]byte) IntCommand {
	return self.Push(false, key, values...)
}

func (self *cache) Pop(beginning bool, key string) BytesCommand {
	var name string
	if beginning {
		name = LPopCommand
	} else {
		name = RPopCommand
	}
	cmdDef := NewCommandDefinition(name, key)
	return self.command(cmdDef)
}

func (self *cache) LPop(key string) BytesCommand {
	return self.Pop(true, key)
}

func (self *cache) RPop(key string) BytesCommand {
	return self.Pop(false, key)
}

func (self *cache) LRange(key string, start int, stop int) BytesSliceCommand {
	cmdDef := NewCommandDefinition(LRangeCommand, key, start, stop)
	return self.command(cmdDef)
}

func (self *cache) LIndex(key string) BytesCommand {
	cmdDef := NewCommandDefinition(LIndexCommand, key)
	return self.command(cmdDef)
}

///////////////////////// hash ////////////////////////
func (self *cache) HKeys(key string) StringSliceCommand {
	cmdDef := NewCommandDefinition(HKeysCommand, key)
	return self.command(cmdDef)
}

func (self *cache) HSet(key string, hashKey []byte, value []byte) BoolCommand {
	cmdDef := NewCommandDefinition(HSetCommand, key, hashKey, value)
	return self.command(cmdDef)
}

func (self *cache) HGet(key string, hashKey []byte) BytesCommand {
	cmdDef := NewCommandDefinition(HGetCommand, key, hashKey)
	return self.command(cmdDef)
}

func (self *cache) HDel(key string, hashKeys ...[]byte) IntCommand {
	args := make([]interface{}, 1+len(hashKeys))
	args[0] = key
	for i, k := range hashKeys {
		args[i+1] = k
	}
	cmdDef := NewCommandDefinition(HDelCommand, args...)
	return self.command(cmdDef)
}

func New(opts *Options) Cache {
	var cache cache

	var auther Auther
	var client Client
	auth := opts.Auth
	if auth == "" {
		auther = new(dummyAuther)
	} else {
		auther = newAuther(NewCommandDefinition(AuthCommand, auth))
	}
	addrs := opts.Addrs
	if len(addrs) > 1 {
		client = newMultiClient(addrs, opts.PoolSize, opts.DialTimeout, auther)
	} else {
		client = newBaseClient(addrs[0], opts.PoolSize, opts.DialTimeout, auther)
	}
	cache.client = client
	return &cache
}
