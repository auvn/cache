package commands

import (
	"errors"
	"log"
	"reflect"

	"github.com/auvn/go.cache/core"
	"github.com/auvn/go.cache/session"
)

const (
	_        int = 1 << iota
	WFlag        // write
	RFlag        // read
	TDFlag       // time dependent
	AuthFlag     // auth required
)

var (
	ErrNonFunc                   = errors.New("non func")
	ErrNilFunc                   = errors.New("nil func")
	ErrCommandExists             = errors.New("func exists")
	ErrArgReflectorExists        = errors.New("argument reflector exists")
	ErrNonErrOut                 = errors.New("out argument type is not error")
	ErrUnknownInArgumentType     = errors.New("cannot find reflector for one of command argument")
	ErrInvalidOutNum             = errors.New("func should return: (interface{}, error)")
	ErrSessionArgPos             = errors.New("session arguments should be first")
	ErrNonValuedVariadicArgument = errors.New("non-valued argument cannot be variadic")

	CoreValueType    = reflect.TypeOf(core.Value{})
	CoreStrValueType = reflect.TypeOf(core.StrValue(""))
	CoreIntValueType = reflect.TypeOf(core.IntValue(0))

	SessionType = reflect.TypeOf((*session.Session)(nil)).Elem()

	ErrorType = reflect.TypeOf((*error)(nil)).Elem()

	Flags = &flags{
		R:  RFlag,
		W:  WFlag,
		RA: RFlag | AuthFlag,
		WA: WFlag | AuthFlag,
		TD: TDFlag,
		A:  AuthFlag,
	}

	DefaultFlag = (AuthFlag)
)

type flags struct {
	R  int
	RA int
	W  int
	WA int
	TD int
	A  int
}

func CheckFlag(flag int, expectedFlag int) bool {
	return flag&expectedFlag != 0
}

type ReflectorInput [2]interface{}

func (self ReflectorInput) Session() session.Session {
	return self[0].(session.Session)
}

func (self ReflectorInput) ArgsIterator() ArgumentsIterator {
	return self[1].(ArgumentsIterator)
}

type ArgumentReflector interface {
	Reflect(in ReflectorInput) (reflect.Value, error)
	Valued() bool
}

type argumentReflectorFn func(in ReflectorInput) (reflect.Value, error)

type valued argumentReflectorFn

func (self valued) Reflect(in ReflectorInput) (reflect.Value, error) {
	return self(in)
}

func (self valued) Valued() bool {
	return true
}

type nonValued argumentReflectorFn

func (self nonValued) Reflect(in ReflectorInput) (reflect.Value, error) {
	return self(in)
}

func (self nonValued) Valued() bool {
	return false
}

type ArgumentsProvider map[reflect.Type]ArgumentReflector

func (self ArgumentsProvider) Get(t reflect.Type) (ArgumentReflector, bool) {
	ar, ok := self[t]
	return ar, ok
}

func (self ArgumentsProvider) GetVariadic(t reflect.Type) (ArgumentReflector, bool) {
	if t.Kind() != reflect.Slice {
		return nil, false
	}
	elemType := t.Elem()
	ar, ok := self.Get(elemType)
	if !ok {
		return nil, ok
	}
	reflectorObj := &variadicReflector{
		t:             t,
		baseReflector: ar,
	}
	return reflectorObj, true
}

func NewArgumentsProvider() ArgumentsProvider {
	return ArgumentsProvider{
		CoreValueType:    valued(ValueReflector),
		CoreStrValueType: valued(StrValueReflector),
		CoreIntValueType: valued(IntValueReflector),
		SessionType:      nonValued(SessionValueReflector),
	}
}

type variadicReflector struct {
	t             reflect.Type
	baseReflector ArgumentReflector
}

func (self *variadicReflector) Reflect(in ReflectorInput) (reflect.Value, error) {
	array := reflect.MakeSlice(self.t, 0, 0)
	var err error
	var value reflect.Value
	for err == nil {
		if value, err = self.baseReflector.Reflect(in); err == nil {
			array = reflect.Append(array, value)
		}
	}
	return array, nil
}

func (self *variadicReflector) Valued() bool {
	return self.baseReflector.Valued()
}

func ValueReflector(in ReflectorInput) (reflect.Value, error) {
	v, err := in.ArgsIterator().Next()
	return reflect.ValueOf(v), err
}

func StrValueReflector(in ReflectorInput) (reflect.Value, error) {
	v, err := in.ArgsIterator().NextStr()
	return reflect.ValueOf(v), err
}

func IntValueReflector(in ReflectorInput) (reflect.Value, error) {
	v, err := in.ArgsIterator().NextInt()
	return reflect.ValueOf(v), err
}

func SessionValueReflector(in ReflectorInput) (reflect.Value, error) {
	return reflect.ValueOf(in.Session()), nil
}

type ReflectRegistryOptions struct {
	AuthEnabled bool
}

type ReflectRegistry struct {
	r    MapRegistry
	args ArgumentsProvider
	opts *ReflectRegistryOptions
}

func (self *ReflectRegistry) updateFlag(flag int) int {
	if !self.opts.AuthEnabled && CheckFlag(flag, AuthFlag) {
		flag = flag ^ AuthFlag
	}
	return flag
}

func (self *ReflectRegistry) Register(name string, fn interface{}, flag int) error {
	if fn == nil {
		return ErrNilFunc
	}

	if _, ok := self.r[name]; ok {
		return ErrCommandExists
	}

	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()
	if fnType.Kind() != reflect.Func {
		return ErrNonFunc
	}
	valuedArgs := 0
	variadic := fnType.IsVariadic()
	numIn := fnType.NumIn()
	args := make([]ArgumentReflector, numIn)

	var reflector ArgumentReflector
	var ok bool
	var variadicArg bool
	for i := 0; i < numIn; i++ {
		arg := fnType.In(i)
		variadicArg = variadic && i+1 == numIn
		if variadicArg {
			reflector, ok = self.args.GetVariadic(arg)
		} else {
			reflector, ok = self.args.Get(arg)
		}
		if !ok {
			return ErrUnknownInArgumentType
		}

		if !reflector.Valued() {
			if variadicArg {
				return ErrNonValuedVariadicArgument
			}
		} else {
			valuedArgs += 1
		}
		args[i] = reflector
	}

	numOut := fnType.NumOut()
	if numOut != 2 {
		return ErrInvalidOutNum
	}

	cmdErrOut := fnType.Out(1)
	if !cmdErrOut.AssignableTo(ErrorType) {
		return ErrNonErrOut
	}

	flag = self.updateFlag(flag)
	self.r[name] = &reflectedCommand{
		fn:         fnValue,
		args:       args,
		valuedArgs: valuedArgs,
		variadic:   variadic,
		flag:       flag,
	}
	return nil
}

func (self *ReflectRegistry) Get(name string) (Command, bool) {
	return self.r.Get(name)
}

func (self *ReflectRegistry) Begin() *ReflectCommandsBuilder {
	return NewReflectCommandsBuilder(self)
}

func NewReflectRegistry(opts *ReflectRegistryOptions) *ReflectRegistry {
	if opts == nil {
		opts = &ReflectRegistryOptions{}
	}
	return &ReflectRegistry{
		r:    MapRegistry{},
		args: NewArgumentsProvider(),
		opts: opts,
	}
}

type reflectCmdDef struct {
	fn   interface{}
	flag int
}

type ReflectCommandsBuilder struct {
	cmdMap map[string]*reflectCmdDef
	r      *ReflectRegistry
}

func (self *ReflectCommandsBuilder) Cmd(name string, fn interface{}, flags ...int) *ReflectCommandsBuilder {
	var flag int
	if len(flags) == 0 {
		flag = DefaultFlag
	} else {
		flag = 0
		for _, f := range flags {
			flag = flag | f
		}
	}
	self.cmdMap[name] = &reflectCmdDef{fn: fn, flag: flag}
	return self
}

func (self *ReflectCommandsBuilder) MustEnd() *ReflectRegistry {
	if err := self.End(); err != nil {
		log.Fatal(err)
	}
	return self.r
}

func (self *ReflectCommandsBuilder) End() error {
	for k, v := range self.cmdMap {
		if err := self.r.Register(k, v.fn, v.flag); err != nil {
			return err
		}
	}
	return nil
}

func NewReflectCommandsBuilder(r *ReflectRegistry) *ReflectCommandsBuilder {
	return &ReflectCommandsBuilder{
		cmdMap: map[string]*reflectCmdDef{},
		r:      r,
	}
}

type reflectedCommand struct {
	fn         reflect.Value
	args       []ArgumentReflector
	valuedArgs int
	variadic   bool
	flag       int
}

func (self *reflectedCommand) IsFlag(flag int) bool {
	return CheckFlag(self.flag, flag)
}

func (self *reflectedCommand) Flag() int {
	return self.flag
}

func (self *reflectedCommand) authed(s session.Session) bool {
	flag := self.IsFlag(AuthFlag)
	if flag && s.Authenticated() {
		return true
	} else if !flag {
		return true
	}
	return false
}

func (self *reflectedCommand) Execute(s session.Session, arguments Arguments) (interface{}, error) {
	if !self.authed(s) {
		return nil, ErrAuthRequired
	}

	var iter ArgumentsIterator
	var err error
	n := self.valuedArgs
	if self.variadic {
		iter, err = arguments.IterAtLeast(n - 1)
	} else {
		iter, err = arguments.IterN(n)
	}
	if err != nil {
		return nil, err
	}

	reflectorInput := ReflectorInput{s, iter}
	in := make([]reflect.Value, len(self.args))
	for i, _ := range in {
		argReflector := self.args[i]
		if a, err := argReflector.Reflect(reflectorInput); err == nil {
			in[i] = a
		} else {
			return nil, err
		}
	}

	var ret []reflect.Value
	if self.variadic {
		ret = self.fn.CallSlice(in)
	} else {
		ret = self.fn.Call(in)
	}

	cmdRet := ret[0].Interface()
	cmdErr, _ := ret[1].Interface().(error)
	return cmdRet, cmdErr
}
