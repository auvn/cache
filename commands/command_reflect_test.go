package commands

import (
	"testing"

	"github.com/auvn/go.cache/core"
	"github.com/auvn/go.cache/session"
)

func TestReflectRegistry_Register(t *testing.T) {
	type fields struct {
		r    MapRegistry
		args ArgumentsProvider
		opts *ReflectRegistryOptions
	}
	type args struct {
		name string
		fn   interface{}
		flag int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "ErrNilFunc",
			fields: fields{},
			args: args{
				name: "cmd",
				fn:   nil,
				flag: 0,
			},
			wantErr: true,
		},
		{
			name: "ErrCommandExists",
			fields: fields{
				r: MapRegistry{
					"cmd": &reflectedCommand{},
				},
			},
			args: args{
				name: "cmd",
				fn:   1,
				flag: 0,
			},
			wantErr: true,
		},
		{
			name:   "ErrNonFunc",
			fields: fields{},
			args: args{
				name: "cmd",
				fn:   1,
				flag: 0,
			},
			wantErr: true,
		},
		{
			name:   "ErrUnknownInArgumentType",
			fields: fields{},
			args: args{
				name: "cmd",
				fn:   func(i int) {},
				flag: 0,
			},
			wantErr: true,
		},
		{
			name: "ErrNonValuedVariadicArgument",
			fields: fields{
				args: NewArgumentsProvider(),
			},
			args: args{
				name: "cmd",
				fn:   func(ss ...session.Session) {},
				flag: 0,
			},
			wantErr: true,
		},
		{
			name: "ErrInvalidOutNum",
			fields: fields{
				args: NewArgumentsProvider(),
			},
			args: args{
				name: "cmd",
				fn:   func(s session.Session) {},
				flag: 0,
			},
			wantErr: true,
		},
		{
			name: "ErrNonErrOut",
			fields: fields{
				args: NewArgumentsProvider(),
			},
			args: args{
				name: "cmd",
				fn:   func(s session.Session) interface{} { return "" },
				flag: 0,
			},
			wantErr: true,
		},
		{
			name: "Valid",
			fields: fields{
				opts: &ReflectRegistryOptions{},
				args: NewArgumentsProvider(),
				r:    MapRegistry{},
			},
			args: args{
				name: "cmd",
				fn:   func(s session.Session, v core.Value, v2s ...core.StrValue) (interface{}, error) { return "", nil },
				flag: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &ReflectRegistry{
				r:    tt.fields.r,
				args: tt.fields.args,
				opts: tt.fields.opts,
			}
			if err := self.Register(tt.args.name, tt.args.fn, tt.args.flag); (err != nil) != tt.wantErr {
				t.Errorf("ReflectRegistry.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReflectRegistry_updateFlag(t *testing.T) {
	type fields struct {
		r    MapRegistry
		args ArgumentsProvider
		opts *ReflectRegistryOptions
	}
	type args struct {
		flag int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "AuthEnabled",
			fields: fields{
				opts: &ReflectRegistryOptions{AuthEnabled: true},
			},
			args: args{flag: AuthFlag},
			want: AuthFlag,
		},
		{
			name: "AuthDisabled",
			fields: fields{
				opts: &ReflectRegistryOptions{AuthEnabled: false},
			},
			args: args{flag: AuthFlag},
			want: 0,
		},
		{
			name: "AuthDisabled",
			fields: fields{
				opts: &ReflectRegistryOptions{AuthEnabled: false},
			},
			args: args{flag: 0},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &ReflectRegistry{
				r:    tt.fields.r,
				args: tt.fields.args,
				opts: tt.fields.opts,
			}
			if got := self.updateFlag(tt.args.flag); got != tt.want {
				t.Errorf("ReflectRegistry.updateFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}
