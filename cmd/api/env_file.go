package api

import (
	"github.com/spf13/afero"
	"io"
	"log"
)

type fileEnvironment struct {
	path *string
	fs   afero.Fs
	pe   PrintEnvironment
	file afero.File
}

// NewWriteFileEnvironment creates new file output environment
func NewWriteFileEnvironment(path *string, fs afero.Fs, defaultpe PrintEnvironment) PrintEnvironment {
	pe := &fileEnvironment{
		path: path,
		fs:   fs,
		pe:   defaultpe,
	}
	return pe
}

func (e *fileEnvironment) create(path *string, fs afero.Fs) error {
	f, err := fs.Create(*path)
	if err != nil {
		return err
	}

	e.pe = NewStringEnvironment(f)
	e.file = f

	return nil
}

func (e *fileEnvironment) NewPrinter() Printer {
	if *e.path == "" {
		return e.pe.NewPrinter()
	}
	err := e.create(e.path, e.fs)
	if err != nil {
		log.Println(err)
		return e.pe.NewPrinter()
	}
	return NewPrinter(e)
}

func (e *fileEnvironment) PrintFunc(w io.Writer, format string, a ...interface{}) {
	e.pe.PrintFunc(w, format, a...)
}

func (e *fileEnvironment) Writer() io.WriteCloser {
	return e.pe.Writer()
}
