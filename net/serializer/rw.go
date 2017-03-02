package serializer

import "io"

type ReadWriter interface {
	Reader
	Writer
}

type readWriter struct {
	Reader
	Writer
}

func NewReadWriter(r io.Reader, w io.Writer) ReadWriter {
	return &readWriter{Reader: NewReader(r), Writer: NewWriter(w)}
}
