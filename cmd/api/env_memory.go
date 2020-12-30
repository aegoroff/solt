package api

import (
	"bytes"
	"io"
)

type memoryEnvironment struct {
	w  *memoryBufferClosable
	pe PrintEnvironment
}

type memoryBufferClosable struct {
	buf *bytes.Buffer
}

func (m *memoryBufferClosable) Write(p []byte) (n int, err error) {
	return m.buf.Write(p)
}

func (m *memoryBufferClosable) Close() error {
	return nil
}

func (m *memoryEnvironment) Writer() io.WriteCloser {
	return m.w
}

func (m *memoryEnvironment) PrintFunc(w io.Writer, format string, a ...interface{}) {
	m.pe.PrintFunc(w, format, a...)
}

func (m *memoryEnvironment) NewPrinter() Printer {
	return m.pe.NewPrinter()
}

func (m *memoryEnvironment) String() string {
	return m.w.buf.String()
}

// NewMemoryEnvironment creates new memory PrintEnvironment implementation
func NewMemoryEnvironment() StringEnvironment {
	bc := &memoryBufferClosable{
		buf: bytes.NewBufferString(""),
	}
	se := NewStringEnvironment(bc)
	return &memoryEnvironment{pe: se, w: bc}
}
