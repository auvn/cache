package commands

import (
	"reflect"
	"testing"
)

func Test_newReflectRegistryOptions(t *testing.T) {
	type args struct {
		opts *RegistryOptions
	}
	tests := []struct {
		name string
		args args
		want *ReflectRegistryOptions
	}{
		{
			name: "AuthEnabled",
			args: args{opts: &RegistryOptions{Auth: "auth"}},
			want: &ReflectRegistryOptions{AuthEnabled: true},
		},
		{
			name: "AuthDisabled",
			args: args{opts: &RegistryOptions{Auth: ""}},
			want: &ReflectRegistryOptions{AuthEnabled: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newReflectRegistryOptions(tt.args.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newReflectRegistryOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
