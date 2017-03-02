package journal

import (
	"errors"
	"log"
	"os"
)

var (
	fileFlag      int         = (os.O_CREATE | os.O_RDWR)
	fileReadFlag  int         = (os.O_CREATE | os.O_RDONLY)
	fileWriteFlag int         = (os.O_CREATE | os.O_WRONLY)
	filePerm      os.FileMode = 0666

	ErrNonEmptyJournal = errors.New("the journal is not empty")
)

func openJournalFile(path string, flag int) (*os.File, error) {
	return os.OpenFile(path, flag, filePerm)
}

func InitFile(path string) (Journal, error) {
	fileReader, err := openJournalFile(path, fileReadFlag)
	if err != nil {
		return nil, err
	}
	fileWriter, err := openJournalFile(path, fileWriteFlag)
	if err != nil {
		return nil, err
	}
	jReader := &reader{r: fileReader}
	jWriter := &writer{ws: fileWriter}
	journal := &fileJournal{r: jReader, w: jWriter, wf: fileWriter}
	return journal, nil
}

func MustInitFile(path string) Journal {
	if j, err := InitFile(path); err != nil {
		log.Fatal("cannot init journal:", err)
	} else {
		return j
	}
	return nil
}

type fileJournal struct {
	readerDone bool
	r          *reader
	w          *writer
	wf         *os.File
}

func (self *fileJournal) NextEntry() ([][]byte, error) {
	if entry, err := self.r.NextEntry(); err != nil && err == ErrEmpty {
		if werr := self.w.seekStart(self.r.pos); werr != nil {
			return emptyEntry, werr
		}
		self.readerDone = true
		return entry, err
	} else {
		return entry, nil
	}
}

func (self *fileJournal) Write(e [][]byte) error {
	if !self.readerDone {
		return ErrNonEmptyJournal
	}
	return self.w.Write(e)
}

func (self *fileJournal) Rollback() error {
	return self.w.Rollback()
}

func (self *fileJournal) Commit() error {
	if err := self.w.Commit(); err != nil {
		return err
	}
	if err := self.wf.Sync(); err != nil {
		return err
	}
	return nil
}
