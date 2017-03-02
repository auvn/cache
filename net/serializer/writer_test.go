package serializer

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

type fakeWriter struct {
	buf bytes.Buffer
}

func (self *fakeWriter) Write([]byte) (int, error) {
	return -1, errors.New("fake")
}

func (self *fakeWriter) String() string {
	return ""
}

type TestWriter interface {
	io.Writer
	String() string
}

func Test_write(t *testing.T) {
	type args struct {
		bs []byte
		w  TestWriter
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name:    "Valid",
			args:    args{bs: []byte("text"), w: &bytes.Buffer{}},
			wantW:   "text",
			wantErr: false,
		},
		{
			name:    "WriteErr",
			args:    args{bs: []byte("text"), w: &fakeWriter{}},
			wantW:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.args.w
			if err := write(w, tt.args.bs); (err != nil) != tt.wantErr {
				t.Errorf("write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("write() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_writeRune(t *testing.T) {
	type args struct {
		r rune
		w TestWriter
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name:    "Valid",
			args:    args{r: '\r', w: &bytes.Buffer{}},
			wantW:   "\r",
			wantErr: false,
		},
		{
			name:    "WriteErr",
			args:    args{r: 'b', w: &fakeWriter{}},
			wantW:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.args.w
			if err := writeRune(w, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("writeRune() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("writeRune() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_writeCRLF(t *testing.T) {
	type args struct {
		w TestWriter
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name:    "Valid",
			args:    args{w: &bytes.Buffer{}},
			wantW:   "\r\n",
			wantErr: false,
		},
		{
			name:    "WriteErr",
			args:    args{w: &fakeWriter{}},
			wantW:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.args.w
			if err := writeCRLF(w); (err != nil) != tt.wantErr {
				t.Errorf("writeCRLF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("writeCRLF() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_writeIntPrefix(t *testing.T) {
	type args struct {
		w TestWriter
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name:    "Valid",
			args:    args{w: &bytes.Buffer{}},
			wantW:   "I",
			wantErr: false,
		},
		{
			name:    "WriteErr",
			args:    args{w: &fakeWriter{}},
			wantW:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.args.w
			if err := writeIntPrefix(w); (err != nil) != tt.wantErr {
				t.Errorf("writeIntPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("writeIntPrefix() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_writeValuePrefix(t *testing.T) {
	type args struct {
		w TestWriter
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name:    "Valid",
			args:    args{w: &bytes.Buffer{}},
			wantW:   "V",
			wantErr: false,
		},
		{
			name:    "WriteErr",
			args:    args{w: &fakeWriter{}},
			wantW:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.args.w
			if err := writeValuePrefix(w); (err != nil) != tt.wantErr {
				t.Errorf("writeValuePrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("writeValuePrefix() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_writeErrPrefix(t *testing.T) {
	type args struct {
		w TestWriter
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name:    "Valid",
			args:    args{w: &bytes.Buffer{}},
			wantW:   "E",
			wantErr: false,
		},
		{
			name:    "WriteErr",
			args:    args{w: &fakeWriter{}},
			wantW:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.args.w
			if err := writeErrPrefix(w); (err != nil) != tt.wantErr {
				t.Errorf("writeErrPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("writeErrPrefix() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_writeArrayPrefix(t *testing.T) {
	type args struct {
		w TestWriter
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name:    "Valid",
			args:    args{w: &bytes.Buffer{}},
			wantW:   "A",
			wantErr: false,
		},
		{
			name:    "WriteErr",
			args:    args{w: &fakeWriter{}},
			wantW:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.args.w
			if err := writeArrayPrefix(w); (err != nil) != tt.wantErr {
				t.Errorf("writeArrayPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("writeArrayPrefix() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_writeNilPrefix(t *testing.T) {
	type args struct {
		w TestWriter
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name:    "Valid",
			args:    args{w: &bytes.Buffer{}},
			wantW:   "N",
			wantErr: false,
		},
		{
			name:    "WriteErr",
			args:    args{w: &fakeWriter{}},
			wantW:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.args.w
			if err := writeNilPrefix(w); (err != nil) != tt.wantErr {
				t.Errorf("writeNilPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("writeNilPrefix() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_writeBoolPrefix(t *testing.T) {
	type args struct {
		w TestWriter
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name:    "Valid",
			args:    args{w: &bytes.Buffer{}},
			wantW:   "B",
			wantErr: false,
		},
		{
			name:    "WriteErr",
			args:    args{w: &fakeWriter{}},
			wantW:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.args.w
			if err := writeBoolPrefix(w); (err != nil) != tt.wantErr {
				t.Errorf("writeBoolPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("writeBoolPrefix() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_writeInt(t *testing.T) {
	type args struct {
		i int64
		w TestWriter
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name:    "Valid",
			args:    args{i: 10, w: &bytes.Buffer{}},
			wantW:   "10",
			wantErr: false,
		},
		{
			name:    "Valid",
			args:    args{i: -10, w: &bytes.Buffer{}},
			wantW:   "-10",
			wantErr: false,
		},
		{
			name:    "WriteErr",
			args:    args{w: &fakeWriter{}},
			wantW:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.args.w
			if err := writeInt(tt.args.w, tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("writeInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("writeInt() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_writeIntValue(t *testing.T) {
	type args struct {
		i int64
		w TestWriter
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name:    "Valid",
			args:    args{i: 10, w: &bytes.Buffer{}},
			wantW:   "I10\r\n",
			wantErr: false,
		},
		{
			name:    "Valid",
			args:    args{i: -10, w: &bytes.Buffer{}},
			wantW:   "I-10\r\n",
			wantErr: false,
		},
		{
			name:    "WriteErr",
			args:    args{i: 10, w: &fakeWriter{}},
			wantW:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.args.w
			if err := writeIntValue(w, tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("writeIntValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("writeIntValue() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_writeBytesValue(t *testing.T) {
	type args struct {
		bs []byte
		w  TestWriter
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name:    "Valid",
			args:    args{bs: []byte("str"), w: &bytes.Buffer{}},
			wantW:   "V3\r\nstr\r\n",
			wantErr: false,
		},

		{
			name:    "Valid",
			args:    args{bs: []byte("\r\n"), w: &bytes.Buffer{}},
			wantW:   "V2\r\n\r\n\r\n",
			wantErr: false,
		},

		{
			name:    "EmptyBytesSlice",
			args:    args{bs: []byte{}, w: &bytes.Buffer{}},
			wantW:   "V0\r\n\r\n",
			wantErr: false,
		},
		{
			name:    "WriteErr",
			args:    args{bs: []byte("value"), w: &fakeWriter{}},
			wantW:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.args.w
			if err := writeBytesValue(w, tt.args.bs); (err != nil) != tt.wantErr {
				t.Errorf("writeBytesValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("writeBytesValue() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_writeErrorValue(t *testing.T) {
	type args struct {
		err error
		w   TestWriter
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name:    "Error",
			args:    args{err: errors.New("err body"), w: &bytes.Buffer{}},
			wantW:   "Eerr body\r\n",
			wantErr: false,
		},
		{
			name:    "WriteErrorErr",
			args:    args{err: errors.New("err body"), w: &fakeWriter{}},
			wantW:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.args.w
			if err := writeErrorValue(w, tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("writeErrorValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("writeErrorValue() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_writeArrayValue(t *testing.T) {
	type args struct {
		arr []interface{}
		w   TestWriter
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name:    "Array",
			args:    args{arr: []interface{}{"value", 10, true, errors.New(""), nil, []byte("str")}, w: &bytes.Buffer{}},
			wantW:   "A6\r\nV5\r\nvalue\r\nI10\r\nB1\r\nE\r\nN\r\nV3\r\nstr\r\n",
			wantErr: false,
		},
		{
			name:    "Err",
			args:    args{arr: []interface{}{"value", 10, true}, w: &fakeWriter{}},
			wantW:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.args.w
			if err := writeArrayValue(w, tt.args.arr); (err != nil) != tt.wantErr {
				t.Errorf("writeArrayValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("writeArrayValue() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_writeNilValue(t *testing.T) {
	type args struct {
		w TestWriter
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name:    "Nil",
			args:    args{&bytes.Buffer{}},
			wantW:   "N\r\n",
			wantErr: false,
		},
		{
			name:    "Err",
			args:    args{&fakeWriter{}},
			wantW:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.args.w
			if err := writeNilValue(w); (err != nil) != tt.wantErr {
				t.Errorf("writeNilValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("writeNilValue() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_writeBoolValue(t *testing.T) {
	type args struct {
		v bool
		w TestWriter
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name:    "Bool",
			args:    args{v: true, w: &bytes.Buffer{}},
			wantW:   "B1\r\n",
			wantErr: false,
		},
		{
			name:    "Err",
			args:    args{v: true, w: &fakeWriter{}},
			wantW:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.args.w
			if err := writeBoolValue(w, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("writeBoolValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("writeBoolValue() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_writeInterface(t *testing.T) {
	type args struct {
		v interface{}
		w TestWriter
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name:    "Nil",
			args:    args{v: nil, w: &bytes.Buffer{}},
			wantW:   "N\r\n",
			wantErr: false,
		},
		{
			name:    "WriteNilErr",
			args:    args{v: nil, w: &fakeWriter{}},
			wantW:   "",
			wantErr: true,
		},
		{
			name:    "Error",
			args:    args{v: errors.New("err body"), w: &bytes.Buffer{}},
			wantW:   "Eerr body\r\n",
			wantErr: false,
		},
		{
			name:    "WriteErrorErr",
			args:    args{v: errors.New("err body"), w: &fakeWriter{}},
			wantW:   "",
			wantErr: true,
		},
		{
			name:    "Bool",
			args:    args{v: true, w: &bytes.Buffer{}},
			wantW:   "B1\r\n",
			wantErr: false,
		},
		{
			name:    "WriteBoolErr",
			args:    args{v: true, w: &fakeWriter{}},
			wantW:   "",
			wantErr: true,
		},
		{
			name:    "Int",
			args:    args{v: 10, w: &bytes.Buffer{}},
			wantW:   "I10\r\n",
			wantErr: false,
		},
		{
			name:    "WriteIntErr",
			args:    args{v: 10, w: &fakeWriter{}},
			wantW:   "",
			wantErr: true,
		},
		{
			name:    "String",
			args:    args{v: "str", w: &bytes.Buffer{}},
			wantW:   "V3\r\nstr\r\n",
			wantErr: false,
		},
		{
			name:    "WriteStringErr",
			args:    args{v: "str", w: &fakeWriter{}},
			wantW:   "",
			wantErr: true,
		},
		{
			name:    "BytesSlice",
			args:    args{v: []byte("slice"), w: &bytes.Buffer{}},
			wantW:   "V5\r\nslice\r\n",
			wantErr: false,
		},
		{
			name:    "WriteBytesSliceErr",
			args:    args{v: []byte("slice"), w: &fakeWriter{}},
			wantW:   "",
			wantErr: true,
		},
		{
			name:    "InterfacesSlice",
			args:    args{v: []interface{}{"value", 10, true, errors.New(""), nil, []byte("str")}, w: &bytes.Buffer{}},
			wantW:   "A6\r\nV5\r\nvalue\r\nI10\r\nB1\r\nE\r\nN\r\nV3\r\nstr\r\n",
			wantErr: false,
		},
		{
			name:    "WriteInterfacesSliceErr",
			args:    args{v: []interface{}{"value", 10, true}, w: &fakeWriter{}},
			wantW:   "",
			wantErr: true,
		},
		{
			name:    "Err",
			args:    args{v: struct{}{}, w: &bytes.Buffer{}},
			wantW:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.args.w
			if err := writeInterface(w, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("writeInterface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("writeInterface() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
