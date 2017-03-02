package journal

const (
	notCreated byte = 0
	initiated       = 1
	progress        = 2
	commited        = 3
	rollback        = 4

	statusOffset int = 1 //byte
)

var (
	emptyEntry     = [][]byte{}
	emptyEntryItem = []byte{}

	Empty Journal = emptyJournal(true)
)

type Reader interface {
	NextEntry() ([][]byte, error)
}

type Writer interface {
	Write([][]byte) error
	Rollback() error
	Commit() error
}

type Journal interface {
	Reader
	Writer
}

type emptyJournal bool

func (self emptyJournal) NextEntry() ([][]byte, error) {
	return emptyEntry, ErrEmpty
}

func (self emptyJournal) Write([][]byte) error {
	return nil
}

func (self emptyJournal) Rollback() error {
	return nil
}

func (self emptyJournal) Commit() error {
	return nil
}
