package serializer

import (
	"bufio"
	"bytes"
	"io"
	"sync"
)

var (
	emptyBytesBuffer    = bytes.NewBuffer([]byte{})
	bufferedReadersPool = &sync.Pool{
		New: func() interface{} {
			return bufio.NewReader(emptyBytesBuffer)
		},
	}
	bufferredWritersPool = &sync.Pool{
		New: func() interface{} {
			return bufio.NewWriter(emptyBytesBuffer)
		},
	}
)

func getBufferedReader(r io.Reader) *bufio.Reader {
	buf := bufferedReadersPool.Get().(*bufio.Reader)
	buf.Reset(r)
	return buf
}

func putBufferedReader(buf *bufio.Reader) {
	buf.Reset(emptyBytesBuffer)
	bufferedReadersPool.Put(buf)
}

func getBufferedWriter(w io.Writer) *bufio.Writer {
	buf := bufferredWritersPool.Get().(*bufio.Writer)
	buf.Reset(w)
	return buf
}

func putBufferedWriter(buf *bufio.Writer) {
	buf.Reset(emptyBytesBuffer)
	bufferredWritersPool.Put(buf)
}
