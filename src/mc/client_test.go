package mc

import (
	"bytes"
	"testing"
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
	c := NewClient(buf, 50, nil)
	return c, buf
}

//////////////////////////////////////////////////////////////

func TestClientCanHandshake(t *testing.T) {
	// client, buf := createClient()
}
