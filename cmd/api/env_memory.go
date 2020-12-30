package api

import (
	"bytes"
	"io"
)

type memoryEnvironment struct {
	buffer *bufferClosable
	pe     PrintEnvironment
}

type bufferClosable struct {
	*bytes.Buffer
}

func (*bufferClosable) Close() error {
	return nil
}

func (m *memoryEnvironment) Writer() io.WriteCloser {
	return m.buffer
}

func (m *memoryEnvironment) PrintFunc(w io.Writer, format string, a ...interface{}) {
	m.pe.PrintFunc(w, format, a...)
}

func (m *memoryEnvironment) NewPrinter() Printer {
	return m.pe.NewPrinter()
}

func (m *memoryEnvironment) String() string {
	return m.buffer.String()
}

// NewMemoryEnvironment creates new memory PrintEnvironment implementation
func NewMemoryEnvironment() StringEnvironment {
	buffer := &bufferClosable{
		bytes.NewBufferString(""),
	}
	se := NewStringEnvironment(buffer)
	return &memoryEnvironment{pe: se, buffer: buffer}
}
