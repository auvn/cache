package serializer

import (
	"errors"
	"reflect"
	"testing"
)

func Test_payload_Array(t *testing.T) {
	type fields struct {
		v interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		want    []Payload
		wantErr bool
	}{
		{
			name:    "Valid",
			fields:  fields{v: make([]Payload, 10)},
			want:    make([]Payload, 10),
			wantErr: false,
		},

		{
			name:    "Invalid",
			fields:  fields{v: ""},
			want:    []Payload{},
			wantErr: true,
		},
		{
			name:    "NilValue",
			fields:  fields{v: nil},
			want:    []Payload{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &payload{
				v: tt.fields.v,
			}
			got, err := self.Array()
			if (err != nil) != tt.wantErr {
				t.Errorf("payload.Array() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("payload.Array() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_payload_Bytes(t *testing.T) {
	type fields struct {
		v interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name:    "Valid",
			fields:  fields{v: make([]byte, 10)},
			want:    make([]byte, 10),
			wantErr: false,
		},
		{
			name:    "Invalid",
			fields:  fields{v: ""},
			want:    []byte{},
			wantErr: true,
		},
		{
			name:    "NilValue",
			fields:  fields{v: nil},
			want:    []byte{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &payload{
				v: tt.fields.v,
			}
			got, err := self.Bytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("payload.Bytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("payload.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_payload_Str(t *testing.T) {
	type fields struct {
		v interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name:    "BytesValue",
			fields:  fields{v: make([]byte, 10)},
			want:    string(make([]byte, 10)),
			wantErr: false,
		},
		{
			name:    "NonBytesValue",
			fields:  fields{v: ""},
			want:    "",
			wantErr: true,
		},
		{
			name:    "NonBytesValue2",
			fields:  fields{v: 1},
			want:    "",
			wantErr: true,
		},
		{
			name:    "NilValue",
			fields:  fields{v: nil},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &payload{
				v: tt.fields.v,
			}
			got, err := self.Str()
			if (err != nil) != tt.wantErr {
				t.Errorf("payload.Str() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("payload.Str() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_payload_Int(t *testing.T) {
	type fields struct {
		v interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{
			name:    "BytesValue",
			fields:  fields{v: []byte("123")},
			want:    123,
			wantErr: false,
		},
		{
			name:    "NonConvertableBytesValue",
			fields:  fields{v: []byte("123b")},
			want:    0,
			wantErr: true,
		},
		{
			name:    "NonBytesValue",
			fields:  fields{v: 1},
			want:    0,
			wantErr: true,
		},
		{
			name:    "NonBytesValue2",
			fields:  fields{v: "123"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "NilValue",
			fields:  fields{v: nil},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &payload{
				v: tt.fields.v,
			}
			got, err := self.Int()
			if (err != nil) != tt.wantErr {
				t.Errorf("payload.Int() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("payload.Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_payload_Bool(t *testing.T) {
	type fields struct {
		v interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{

		{
			name:    "BytesValueFalse",
			fields:  fields{v: []byte("0")},
			want:    false,
			wantErr: false,
		},
		{
			name:    "BytesValueTrue",
			fields:  fields{v: []byte("1")},
			want:    true,
			wantErr: false,
		},
		{
			name:    "BytesValueTrue",
			fields:  fields{v: []byte("b")},
			want:    true,
			wantErr: false,
		},
		{
			name:    "BytesValueBigLen",
			fields:  fields{v: []byte("12")},
			want:    false,
			wantErr: true,
		},
		{
			name:    "NonBytesValue",
			fields:  fields{v: "123"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "NonBytesValue",
			fields:  fields{v: true},
			want:    false,
			wantErr: true,
		},
		{
			name:    "NilValue",
			fields:  fields{v: nil},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &payload{
				v: tt.fields.v,
			}
			got, err := self.Bool()
			if (err != nil) != tt.wantErr {
				t.Errorf("payload.Bool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("payload.Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_payload_Err(t *testing.T) {
	type fields struct {
		v interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   error
	}{
		{
			name:   "ErrorValue",
			fields: fields{v: errors.New("test")},
			want:   errors.New("test"),
		},
		{
			name:   "NonErrorValue",
			fields: fields{v: "test"},
			want:   ErrPayloadNonError,
		},
		{
			name:   "NilValue",
			fields: fields{v: nil},
			want:   ErrPayloadNonError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &payload{
				v: tt.fields.v,
			}
			// if err := self.Err(); (err != nil) != tt.wantErr {
			//	t.Errorf("payload.Err() error = %v, wantErr %v", err, tt.wantErr)
			// }
			if got := self.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("payload.Err() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_payload_IsNil(t *testing.T) {
	type fields struct {
		v interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "NilValue",
			fields: fields{v: nil},
			want:   true,
		},
		{
			name:   "NonNilValue",
			fields: fields{v: ""},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &payload{
				v: tt.fields.v,
			}
			if got := self.IsNil(); got != tt.want {
				t.Errorf("payload.IsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_payload_IsArray(t *testing.T) {
	type fields struct {
		v interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "Valid",
			fields: fields{v: make([]Payload, 10)},
			want:   true,
		},

		{
			name:   "Invalid",
			fields: fields{v: ""},
			want:   false,
		},
		{
			name:   "NilValue",
			fields: fields{v: nil},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &payload{
				v: tt.fields.v,
			}
			if got := self.IsArray(); got != tt.want {
				t.Errorf("payload.IsArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_payload_IsErr(t *testing.T) {
	type fields struct {
		v interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "ErrorValue",
			fields: fields{v: errors.New("test")},
			want:   true,
		},
		{
			name:   "NonErrorValue",
			fields: fields{v: "test"},
			want:   false,
		},
		{
			name:   "NilValue",
			fields: fields{v: nil},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &payload{
				v: tt.fields.v,
			}
			if got := self.IsErr(); got != tt.want {
				t.Errorf("payload.IsErr() = %v, want %v", got, tt.want)
			}
		})
	}
}
