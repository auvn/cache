package commands

type RegistryOptions struct {
	Auth string
}

func newReflectRegistryOptions(opts *RegistryOptions) *ReflectRegistryOptions {
	registryOptions := new(ReflectRegistryOptions)
	if opts.Auth != "" {
		registryOptions.AuthEnabled = true
	}
	return registryOptions
}

func InitReflectRegistry(opts *RegistryOptions) *ReflectRegistry {
	securityCommand := NewSecurityCommand(opts.Auth)
	storageCommand := NewStorageCommand()
	stringCommand := NewStringCommand()
	listCommand := NewListCommand()
	hashCommand := NewHashCommand()

	registryOptions := newReflectRegistryOptions(opts)
	return NewReflectRegistry(registryOptions).
		Begin().
		//auth
		Cmd("AUTH", securityCommand.Auth, Flags.R).
		//common
		Cmd("KEYS", storageCommand.Keys, Flags.RA).
		Cmd("EXPIRE", storageCommand.Expire, Flags.WA).
		Cmd("DEL", storageCommand.Del, Flags.WA).
		Cmd("TTL", storageCommand.TTL, Flags.RA).
		//string
		Cmd("SET", stringCommand.Set, Flags.WA).
		Cmd("GET", stringCommand.Get, Flags.RA).
		//list
		Cmd("LPUSH", listCommand.LPush, Flags.WA).
		Cmd("RPUSH", listCommand.RPush, Flags.WA).
		Cmd("LPOP", listCommand.LPop, Flags.WA).
		Cmd("RPOP", listCommand.RPop, Flags.WA).
		Cmd("LRANGE", listCommand.LRange, Flags.RA).
		Cmd("LINDEX", listCommand.LIndex, Flags.RA).
		//hash
		Cmd("HSET", hashCommand.Set, Flags.WA).
		Cmd("HGET", hashCommand.Get, Flags.RA).
		Cmd("HDEL", hashCommand.Del, Flags.WA).
		Cmd("HKEYS", hashCommand.Keys, Flags.RA).
		MustEnd()
}
