package journal

import (
	"encoding/binary"
	"io"
)

var (
	_ Writer = (*writer)(nil)
)

type writer struct {
	pos       int64
	prevPos   int64
	statusPos int64
	endPos    int64
	ws        io.WriteSeeker
}

func (self *writer) seekStart(offset int64) error {
	if offset < 0 {
		offset = 0
	}
	newOffset, err := self.ws.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}
	self.pos = newOffset
	return nil
}

func (self *writer) writeBytes(bs []byte) error {
	n, err := self.ws.Write(bs)
	if err != nil {
		return err
	}
	self.pos += int64(n)
	return nil
}

func (self *writer) writeStatus(status byte) error {
	if err := self.writeBytes([]byte{status}); err != nil {
		return err
	}
	return nil
}

func (self *writer) updateStatus(status byte) error {
	if err := self.seekStart(self.statusPos); err != nil {
		return err
	}
	if err := self.writeStatus(status); err != nil {
		return err
	}
	return nil
}

func (self *writer) writeLen(llen int64) error {
	llenValue := make([]byte, 8)
	binary.BigEndian.PutUint64(llenValue, uint64(llen))
	if err := self.writeBytes(llenValue); err != nil {
		return err
	}
	return nil
}

func (self *writer) writeEntry(entry [][]byte) error {
	return nil
}

func (self *writer) Write(entry [][]byte) error {
	entrySize := len(entry)
	if entrySize <= 0 {
		return nil
	}

	self.statusPos = self.pos
	if err := self.writeStatus(initiated); err != nil {
		return err
	}
	if err := self.writeLen(int64(entrySize)); err != nil {
		return err
	}

	var itemSize int

	for _, item := range entry {
		itemSize = len(item)
		if err := self.writeLen(int64(itemSize)); err != nil {
			return err
		}
		if err := self.writeBytes(item); err != nil {
			return err
		}
	}

	// writing status for the next journal item
	if err := self.writeStatus(notCreated); err != nil {
		return err
	}
	self.endPos = self.pos - int64(statusOffset)

	if err := self.updateStatus(progress); err != nil {
		return err
	}
	return nil
}

func (self *writer) Rollback() error {
	if err := self.updateStatus(rollback); err != nil {
		return err
	}
	if err := self.seekStart(self.statusPos); err != nil {
		return err
	}
	return nil
}

//TODO: add check for status, e.g. if progress then commit, else error
func (self *writer) Commit() error {
	if err := self.updateStatus(commited); err != nil {
		return err
	}
	if err := self.seekStart(self.endPos); err != nil {
		return err
	}
	return nil
}

func NewWriter(ws io.WriteSeeker) *writer {
	return &writer{ws: ws}
}
