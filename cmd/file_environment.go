package cmd

import (
	"github.com/aegoroff/dirstat/scan"
	"github.com/spf13/afero"
	"io"
	"log"
	"solt/cmd/api"
)

type fileEnvironment struct {
	path *string
	fs   afero.Fs
	pe   api.PrintEnvironment
	file afero.File
}

func newWriteFileEnvironment(path *string, fs afero.Fs, defaultpe api.PrintEnvironment) *fileEnvironment {
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

	e.pe = api.NewStringEnvironment(f)
	e.file = f

	return nil
}

func (e *fileEnvironment) close() {
	scan.Close(e.file)
}

func (e *fileEnvironment) NewPrinter() api.Printer {
	if *e.path == "" {
		return e.pe.NewPrinter()
	}
	err := e.create(e.path, e.fs)
	if err != nil {
		log.Println(err)
		return e.pe.NewPrinter()
	}
	return api.NewPrinter(e)
}

func (e *fileEnvironment) PrintFunc(w io.Writer, format string, a ...interface{}) {
	e.pe.PrintFunc(w, format, a...)
}

func (e *fileEnvironment) Writer() io.Writer {
	return e.pe.Writer()
}
