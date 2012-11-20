package mc

import (
	"bytes"
)

type ClosableBuffer struct {
	bytes.Buffer
	WasClosed bool
}

func newClosableBuffer() *ClosableBuffer {
	return &ClosableBuffer{Buffer: *bytes.NewBuffer([]byte{})}
}

func (b *ClosableBuffer) Close() error {
	b.WasClosed = true
	return nil
}

//////////////////////////////////////////////////////////////

func createClient() (*Client, *ClosableBuffer) {
	buf := newClosableBuffer()
	c := NewClient(buf)
	return c, buf
}
