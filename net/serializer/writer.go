package serializer

import (
	"io"
	"reflect"
	"strconv"
)

var (
	ErrCannotWrite = NewError("cannot write passed value")
)

// just to ignore n return value
func write(w io.Writer, bs []byte) error {
	if _, err := w.Write(bs); err != nil {
		return err
	}
	return nil
}

func writeRune(w io.Writer, r rune) error {
	return write(w, []byte(string(r)))
}

func writeCRLF(w io.Writer) error {
	return write(w, CRLF)
}

func writeIntPrefix(w io.Writer) error {
	return writeRune(w, IntPrefix)
}

func writeValuePrefix(w io.Writer) error {
	return writeRune(w, ValuePrefix)
}

func writeErrPrefix(w io.Writer) error {
	return writeRune(w, ErrPrefix)
}

func writeArrayPrefix(w io.Writer) error {
	return writeRune(w, ArrayPrefix)
}

func writeNilPrefix(w io.Writer) error {
	return writeRune(w, NilPrefix)
}

func writeBoolPrefix(w io.Writer) error {
	return writeRune(w, BoolPrefix)
}

func writeInt(w io.Writer, i int64) error {
	return write(w, []byte(strconv.FormatInt(i, 10)))
}

func writeIntValue(w io.Writer, i int64) error {
	if err := writeIntPrefix(w); err != nil {
		return err
	}
	if err := writeInt(w, i); err != nil {
		return err
	}
	return writeCRLF(w)
}

func writeBytesValue(w io.Writer, bs []byte) error {
	bslen := int64(len(bs))
	if err := writeValuePrefix(w); err != nil {
		return err
	}
	if err := writeInt(w, bslen); err != nil {
		return err
	}
	if err := writeCRLF(w); err != nil {
		return err
	}
	if err := write(w, bs); err != nil {
		return err
	}
	return writeCRLF(w)
}

func writeErrorValue(w io.Writer, err error) error {
	if err := writeErrPrefix(w); err != nil {
		return err
	}

	if err := write(w, []byte(err.Error())); err != nil {
		return err
	}

	return writeCRLF(w)
}

func writeArrayValue(w io.Writer, arr []interface{}) error {
	if err := writeArrayPrefix(w); err != nil {
		return err
	}
	if err := writeInt(w, int64(len(arr))); err != nil {
		return err
	}
	if err := writeCRLF(w); err != nil {
		return err
	}

	for _, i := range arr {
		if err := writeInterface(w, i); err != nil {
			return err
		}
	}
	return nil
}

func writeNilValue(w io.Writer) error {
	if err := writeNilPrefix(w); err != nil {
		return err
	}
	return writeCRLF(w)
}

func writeBoolValue(w io.Writer, v bool) error {
	if err := writeBoolPrefix(w); err != nil {
		return err
	}
	b := '0'
	if v {
		b = '1'
	}
	if err := writeRune(w, b); err != nil {
		return err
	}

	return writeCRLF(w)
}

func writeInterface(w io.Writer, v interface{}) error {
	if v == nil {
		return writeNilValue(w)
	}

	switch v.(type) {
	case error:
		return writeErrorValue(w, v.(error))
	}

	vValue := reflect.ValueOf(v)
	vType := vValue.Type()
	vKind := vType.Kind()

	switch vKind {
	case reflect.Bool:
		return writeBoolValue(w, vValue.Bool())
	case reflect.Int:
		return writeIntValue(w, vValue.Int())
	case reflect.String:
		return writeBytesValue(w, []byte(vValue.String()))
	case reflect.Slice:
		if vType.Elem().Kind() == reflect.Uint8 {
			return writeBytesValue(w, vValue.Bytes())
		}
		sliceLen := vValue.Len()
		array := make([]interface{}, sliceLen)
		for i, _ := range array {
			array[i] = vValue.Index(i).Interface()
		}
		return writeArrayValue(w, array)
	}
	return ErrCannotWrite
}

func Write(w io.Writer, v interface{}) error {
	return writeInterface(w, v)
}

type Writer interface {
	Write(interface{}) error
}

type writer struct {
	w io.Writer
}

func (self *writer) Write(v interface{}) error {
	buffer := getBufferedWriter(self.w)
	defer putBufferedWriter(buffer)
	if err := writeInterface(buffer, v); err != nil {
		return err
	}
	return buffer.Flush()
}

func NewWriter(w io.Writer) Writer {
	return &writer{w}
}
