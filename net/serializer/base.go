package serializer

const (
	ArrayPrefix = 'A'
	ValuePrefix = 'V'
	IntPrefix   = 'I'
	ErrPrefix   = 'E'
	NilPrefix   = 'N'
	BoolPrefix  = 'B'

	CR = '\r'
	LF = '\n'

	MaxValueSize = 256 * 1024 * 1024 //256MB
)

var CRLF = []byte("\r\n")
