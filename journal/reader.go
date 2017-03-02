package journal

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	ErrEmpty = errors.New("empty journal")
	ErrEOF   = errors.New("end of journal")

	_ Reader = (*reader)(nil)
)

type reader struct {
	pos int64
	r   io.Reader
}

func (self *reader) readBytes(bs []byte) error {
	_, err := io.ReadFull(self.r, bs)
	if err != nil {
		return err
	}
	self.pos += int64(len(bs))
	return nil
}

func (self *reader) readStatus() (byte, error) {
	statusValue := make([]byte, 1)
	if err := self.readBytes(statusValue); err != nil {
		if err == io.EOF {
			return 0, ErrEOF
		} else {
			return 0, err
		}
	}
	return statusValue[0], nil
}

func (self *reader) readLen() (int64, error) {
	llenValue := make([]byte, 8)
	if err := self.readBytes(llenValue); err != nil {
		return 0, err
	}
	val := binary.BigEndian.Uint64(llenValue)
	return int64(val), nil
}

func (self *reader) readEntryItem() ([]byte, error) {
	itemLen, err := self.readLen()
	if err != nil {
		return emptyEntryItem, err
	}
	item := make([]byte, itemLen)
	if err := self.readBytes(item); err != nil {
		return emptyEntryItem, err
	}
	return item, nil
}

func (self *reader) readEntry() ([][]byte, error) {
	llen, err := self.readLen()
	if err != nil {
		return emptyEntry, err
	}
	entry := make([][]byte, llen)
	for i, _ := range entry {
		if item, err := self.readEntryItem(); err == nil {
			entry[i] = item
		} else {
			return emptyEntry, err
		}
	}
	return entry, nil
}

func (self *reader) NextEntry() ([][]byte, error) {
	status, err := self.readStatus()
	if err != nil && err != ErrEOF {
		return emptyEntry, err
	}

	switch status {
	case commited:
		return self.readEntry()
	default:
		self.pos -= int64(statusOffset)
		return emptyEntry, ErrEmpty
	}
}

func (self *reader) Tell() int64 {
	return self.pos
}

func NewReader(r io.Reader) *reader {
	return &reader{r: r}
}
