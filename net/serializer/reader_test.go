package serializer

import (
	"bufio"
	"bytes"
	"errors"
	"reflect"
	"testing"
)

func NewBuffer(str string) *bufio.Reader {
	return bufio.NewReader(bytes.NewBufferString(str))
}

func NewZeroBuffer() *bufio.Reader {
	return bufio.NewReader(bytes.NewBufferString(""))
}

func Test_readPrefix(t *testing.T) {
	type args struct {
		buf *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    rune
		wantErr bool
	}{
		{
			name:    "Valid",
			args:    args{buf: NewBuffer("Ptext")},
			want:    'P',
			wantErr: false,
		},
		{
			name:    "Valid",
			args:    args{buf: NewBuffer(":text")},
			want:    ':',
			wantErr: false,
		},
		{
			name:    "EmptyBuffer",
			args:    args{buf: NewBuffer("")},
			want:    0,
			wantErr: true,
		},
		{
			name:    "SkipCRLF",
			args:    args{buf: NewBuffer("\r\nD")},
			want:    'D',
			wantErr: false,
		},
		{
			name:    "SkipCRLFAndEmptyBuffer",
			args:    args{buf: NewBuffer("\r\n")},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readPrefix(tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("readPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readLine(t *testing.T) {
	type args struct {
		buf *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "Valid",
			args:    args{buf: NewBuffer("line\r\n")},
			want:    []byte("line"),
			wantErr: false,
		},
		{
			name:    "Valid",
			args:    args{buf: NewBuffer("\r\nline\r\n")},
			want:    []byte(""),
			wantErr: false,
		},
		{
			name:    "Valid",
			args:    args{buf: NewBuffer("\rline\r\n")},
			want:    []byte("\rline"),
			wantErr: false,
		},
		{
			name:    "NoCR",
			args:    args{buf: NewBuffer("line\n")},
			want:    []byte{},
			wantErr: true,
		},
		{
			name:    "NoLF",
			args:    args{buf: NewBuffer("line\r")},
			want:    []byte{},
			wantErr: true,
		},
		{
			name:    "NoCRLF",
			args:    args{buf: NewBuffer("line")},
			want:    []byte{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readLine(tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("readLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readInt(t *testing.T) {
	type args struct {
		buffer *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "Valid",
			args:    args{NewBuffer("123\r\n")},
			want:    123,
			wantErr: false,
		},
		{
			name:    "Invalid",
			args:    args{NewBuffer("b123\r\n")},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Invalid",
			args:    args{NewBuffer("b123")},
			want:    0,
			wantErr: true,
		},
		{
			name:    "IntOverflow",
			args:    args{NewBuffer("9223372036854775808\r\n")},
			want:    0,
			wantErr: true,
		},
		{
			name:    "IntOverflow",
			args:    args{NewBuffer("-9223372036854775809\r\n")},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readInt(tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("readInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readValueSize(t *testing.T) {
	type args struct {
		buffer *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "Valid",
			args:    args{NewBuffer("12\r\n")},
			want:    12,
			wantErr: false,
		},
		{
			name:    "NonInt",
			args:    args{NewBuffer("aa")},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readValueSize(tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("readValueSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readValueSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readArrayLen(t *testing.T) {
	type args struct {
		buffer *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "Valid",
			args:    args{NewBuffer("12\r\n")},
			want:    12,
			wantErr: false,
		},
		{
			name:    "NonInt",
			args:    args{NewBuffer("aa")},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readArrayLen(tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("readArrayLen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readArrayLen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readArrayPayload(t *testing.T) {
	type args struct {
		buffer *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    Payload
		wantErr bool
	}{
		{
			name: "TwoInts",
			args: args{NewBuffer("2\r\nI256\r\nI256\r\n")},
			want: &payload{v: []Payload{
				&payload{v: []byte("256")},
				&payload{v: []byte("256")},
			}},
			wantErr: false,
		},

		{
			name: "ValueInt",
			args: args{NewBuffer("2\r\nV2\r\nV2\r\nI256\r\n")},
			want: &payload{v: []Payload{
				&payload{v: []byte("V2")},
				&payload{v: []byte("256")},
			}},
			wantErr: false,
		},
		{
			name: "CRLF",
			args: args{NewBuffer("2\r\nV2\r\n\r\n\r\nV2\r\n\r\n\r\n")},
			want: &payload{v: []Payload{
				&payload{v: []byte("\r\n")},
				&payload{v: []byte("\r\n")},
			}},
			wantErr: false,
		},
		{
			name:    "InvalidLen",
			args:    args{NewBuffer("b\r\n")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "NegativeLen",
			args:    args{NewBuffer("-2\r\n")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "InvalidArrayItem",
			args:    args{NewBuffer("2\r\nI5\r\n6")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readArrayPayload(tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("readArrayPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readArrayPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readValuePayload(t *testing.T) {
	type args struct {
		buffer *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    Payload
		wantErr bool
	}{
		{
			name:    "Valid",
			args:    args{NewBuffer("5\r\nthree\r\n")},
			want:    &payload{v: []byte("three")},
			wantErr: false,
		},
		{
			name:    "NoCRLF",
			args:    args{NewBuffer("5\r\nthree")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "CRLF",
			args:    args{NewBuffer("2\r\n\r\n\r\n")},
			want:    &payload{v: []byte("\r\n")},
			wantErr: false,
		},
		{
			name:    "TooLarge",
			args:    args{NewBuffer("262145\r\n")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "TooShort",
			args:    args{NewBuffer("262145\r\nvalue\r\n")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "EmptyValue",
			args:    args{NewBuffer("0\r\n\r\n")},
			want:    &payload{v: []byte{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readValuePayload(tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("readValuePayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readValuePayload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readIntPayload(t *testing.T) {
	type args struct {
		buffer *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    Payload
		wantErr bool
	}{
		{
			name:    "Valid",
			args:    args{NewBuffer("2\r\n")},
			want:    &payload{v: []byte("2")},
			wantErr: false,
		},
		{
			name:    "EmptySlice",
			args:    args{NewBuffer("\r\n")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "EmptyBuffer",
			args:    args{NewZeroBuffer()},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readIntPayload(tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("readIntPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readIntPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readPayloadByPrefix(t *testing.T) {
	type args struct {
		buffer *bufio.Reader
		prefix rune
	}
	tests := []struct {
		name    string
		args    args
		want    Payload
		wantErr bool
	}{
		{
			name:    "InvalidBody",
			args:    args{NewBuffer(".\r\n"), '.'},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Value",
			args:    args{NewBuffer("2\r\n00\r\n"), 'V'},
			want:    &payload{v: []byte("00")},
			wantErr: false,
		},
		{
			name:    "ValueErr",
			args:    args{NewBuffer("\r"), 'V'},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Int",
			args:    args{NewBuffer("123\r\n"), 'I'},
			want:    &payload{v: []byte("123")},
			wantErr: false,
		},
		{
			name:    "IntErr",
			args:    args{NewBuffer("\r"), 'I'},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Array",
			args:    args{NewBuffer("0\r\n"), 'A'},
			want:    &payload{v: []Payload{}},
			wantErr: false,
		},
		{
			name:    "ArrayErr",
			args:    args{NewBuffer("\r"), 'A'},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Bool",
			args:    args{NewBuffer("1\r\n"), 'B'},
			want:    &payload{v: []byte("1")},
			wantErr: false,
		},
		{
			name:    "BoolErr",
			args:    args{NewBuffer("\r"), 'B'},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Err",
			args:    args{NewBuffer("err\r\n"), 'E'},
			want:    &payload{v: errors.New("err")},
			wantErr: false,
		},
		{
			name:    "ErrNil",
			args:    args{NewBuffer("\r"), 'E'},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Nil",
			args:    args{NewBuffer("\r\n"), 'N'},
			want:    &payload{v: nil},
			wantErr: false,
		},
		{
			name:    "NilErr",
			args:    args{NewBuffer("\r"), 'N'},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readPayloadByPrefix(tt.args.buffer, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("readPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lookupPrefix(t *testing.T) {
	type args struct {
		buffer *bufio.Reader
		prefix rune
	}
	tests := []struct {
		name    string
		args    args
		want    rune
		wantErr bool
	}{
		{
			name:    "Valid",
			args:    args{NewBuffer("p"), 'p'},
			want:    'p',
			wantErr: false,
		},
		{
			name:    "Invalid",
			args:    args{NewBuffer("p"), 'P'},
			want:    'p',
			wantErr: true,
		},
		{
			name:    "ReadPrefixError",
			args:    args{NewBuffer(""), 'P'},
			want:    0,
			wantErr: true,
		},
		{
			name:    "SkipCRLF",
			args:    args{NewBuffer("\r\nP"), 'P'},
			want:    'P',
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lookupPrefix(tt.args.buffer, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("lookupPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("lookupPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readPayload(t *testing.T) {
	type args struct {
		buffer *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    Payload
		wantErr bool
	}{
		{
			name:    "ReadPrefixError",
			args:    args{NewBuffer("")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Valid",
			args:    args{NewBuffer("I12\r\n")},
			want:    &payload{v: []byte("12")},
			wantErr: false,
		},
		{
			name:    "Err",
			args:    args{NewBuffer("I\r\n")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readPayload(tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("readPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readErrPayload(t *testing.T) {
	type args struct {
		buffer *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    Payload
		wantErr bool
	}{
		{
			name:    "ReadLineErr",
			args:    args{NewBuffer("error\r")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Valid",
			args:    args{NewBuffer("error\r\n")},
			want:    &payload{v: errors.New("error")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readErrPayload(tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("readErrPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readErrPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readNilPayload(t *testing.T) {
	type args struct {
		buffer *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    Payload
		wantErr bool
	}{
		{
			name:    "ReadLineErr",
			args:    args{NewBuffer("N\r")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "NonEmptyLine",
			args:    args{NewBuffer("N\r\n")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Valid",
			args:    args{NewBuffer("\r\n")},
			want:    &payload{v: nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readNilPayload(tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("readNilPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readNilPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readBoolPayload(t *testing.T) {
	type args struct {
		buffer *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    Payload
		wantErr bool
	}{
		{
			name:    "ReadLineErr",
			args:    args{NewBuffer("B\r")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "TwoBytesInLine",
			args:    args{NewBuffer("11\r\n")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Valid",
			args:    args{NewBuffer("1\r\n")},
			want:    &payload{v: []byte("1")},
			wantErr: false,
		},
		{
			name:    "Valid",
			args:    args{NewBuffer("0\r\n")},
			want:    &payload{v: []byte("0")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readBoolPayload(tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("readBoolPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readBoolPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readPositiveInt(t *testing.T) {
	type args struct {
		buffer *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{

		{
			name:    "Valid",
			args:    args{NewBuffer("0\r\n")},
			want:    0,
			wantErr: false,
		},
		{
			name:    "Valid",
			args:    args{NewBuffer("10\r\n")},
			want:    10,
			wantErr: false,
		},
		{
			name:    "ReadIntErr",
			args:    args{NewBuffer("\r\n")},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Negative",
			args:    args{NewBuffer("-1\r\n")},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readPositiveInt(tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("readPositiveInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readPositiveInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
