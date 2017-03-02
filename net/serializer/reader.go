package serializer

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

var (
	ErrArrayPrefixExpected = newPrefixError(ArrayPrefix)
	ErrValuePrefixExpected = newPrefixError(ValuePrefix)

	ErrNotExpectedPrefix     Error = NewError("prefix is not expected")
	ErrInvalidBody           Error = NewError("invalid body")
	ErrIntExpected           Error = NewError("int expected")
	ErrValueTooLarge         Error = NewError("too large")
	ErrNonEmptySliceExpected Error = NewError("non-empty slice expected")
	ErrNonPositiveInt        Error = NewError("positive int expected")
)

func readPrefix(buf *bufio.Reader) (rune, error) {
	for {
		prefix, _, err := buf.ReadRune()
		if err != nil {
			return 0, err
		}
		if prefix != CR && prefix != LF {
			return prefix, nil
		}
	}
}

func readLine(buf *bufio.Reader) ([]byte, error) {
	bs, err := buf.ReadBytes(LF)
	if err != nil {
		return []byte{}, err
	}
	if bytes.HasSuffix(bs, CRLF) {
		return bs[:len(bs)-2], nil
	}
	return []byte{}, io.ErrUnexpectedEOF
}

func readInt(buffer *bufio.Reader) (int, error) {
	if p, err := readIntPayload(buffer); err == nil {
		if i, err := p.Int(); err == nil {
			return i, nil
		}
	}
	return 0, ErrIntExpected
}

func readPositiveInt(buffer *bufio.Reader) (int, error) {
	if i, err := readInt(buffer); err != nil {
		return 0, err
	} else {
		if i < 0 {
			return 0, ErrNonPositiveInt
		} else {
			return i, nil
		}
	}
}

func readValueSize(buffer *bufio.Reader) (int, error) {
	return readPositiveInt(buffer)
}

func readArrayLen(buffer *bufio.Reader) (int, error) {
	return readPositiveInt(buffer)
}

func readValuePayload(buffer *bufio.Reader) (Payload, error) {
	size, err := readValueSize(buffer)
	if err != nil {
		return nil, err
	}

	if size > MaxValueSize {
		return nil, ErrValueTooLarge
	}

	// expecting CRLF after the value
	// e.g. V2\r\nAB\r\n
	sizeWithCRLF := size + len(CRLF)
	valueBuf := make([]byte, sizeWithCRLF)
	if _, err := io.ReadAtLeast(buffer, valueBuf, sizeWithCRLF); err != nil {
		return nil, err
	}

	if !bytes.HasPrefix(valueBuf[size:], CRLF) {
		return nil, ErrInvalidBody
	}

	return &payload{v: valueBuf[:size]}, nil
}

func readIntPayload(buffer *bufio.Reader) (Payload, error) {
	bs, err := readLine(buffer)
	if err != nil {
		return nil, err
	}

	if len(bs) == 0 {
		return nil, ErrInvalidBody
	}

	return &payload{v: bs}, nil
}

func readArrayPayload(buffer *bufio.Reader) (Payload, error) {
	len, err := readArrayLen(buffer)
	if err != nil {
		return nil, err
	}

	array := make([]Payload, len)
	for i, _ := range array {
		value, err := readPayload(buffer)
		if err != nil {
			return nil, err
		}
		array[i] = value
	}
	return &payload{v: array}, nil
}

func readErrPayload(buffer *bufio.Reader) (Payload, error) {
	line, err := readLine(buffer)
	if err != nil {
		return nil, err
	}
	errStr := string(line)
	return &payload{v: errors.New(errStr)}, nil
}

func readNilPayload(buffer *bufio.Reader) (Payload, error) {
	line, err := readLine(buffer)
	if err != nil {
		return nil, err
	}
	if len(line) != 0 {
		return nil, ErrInvalidBody
	}
	return &payload{v: nil}, nil
}

func readBoolPayload(buffer *bufio.Reader) (Payload, error) {
	line, err := readLine(buffer)
	if err != nil {
		return nil, err
	}
	if len(line) != 1 {
		return nil, ErrInvalidBody
	}
	return &payload{v: line}, nil
}

func readPayloadByPrefix(buffer *bufio.Reader, prefix rune) (Payload, error) {
	switch prefix {
	case ValuePrefix:
		return readValuePayload(buffer)
	case IntPrefix:
		return readIntPayload(buffer)
	case ArrayPrefix:
		return readArrayPayload(buffer)
	case BoolPrefix:
		return readBoolPayload(buffer)
	case ErrPrefix:
		return readErrPayload(buffer)
	case NilPrefix:
		return readNilPayload(buffer)
	}
	return nil, ErrInvalidBody
}

func readPayload(buffer *bufio.Reader) (Payload, error) {
	prefix, err := readPrefix(buffer)
	if err != nil {
		return nil, err
	}
	return readPayloadByPrefix(buffer, prefix)
}

func lookupPrefix(buffer *bufio.Reader, prefix rune) (rune, error) {
	p, err := readPrefix(buffer)
	if err != nil {
		return p, err
	}
	if p == prefix {
		return p, nil
	} else {
		return p, newPrefixError(prefix).WithGot(p)
	}
}

func Read(r io.Reader) (Payload, error) {
	buffer := bufio.NewReader(r)
	return readPayload(buffer)
}

type Reader interface {
	ReadArray() (Payload, error)
	Read() (Payload, error)
}

type reader struct {
	r io.Reader
}

//todo: pool for bufffers

func (self *reader) Read() (Payload, error) {
	buffer := getBufferedReader(self.r)
	defer putBufferedReader(buffer)

	return readPayload(buffer)
}

func (self *reader) ReadArray() (Payload, error) {
	buffer := getBufferedReader(self.r)
	defer putBufferedReader(buffer)

	prefix, err := lookupPrefix(buffer, ArrayPrefix)
	if err != nil {
		return nil, err
	}
	p, err := readPayloadByPrefix(buffer, prefix)
	if err != nil {
		return nil, err
	}
	if p.IsArray() {
		return p, nil
	}
	return nil, ErrInvalidBody
}

func NewReader(r io.Reader) Reader {
	return &reader{r: r}
}
