package sys

import (
	"fmt"
	"github.com/spf13/afero"
	"io"
	"log"
	"os"
)

// Filer defines module that works with files
type Filer interface {
	// CheckExistence validates files passed to be present in file system
	// The list of non exist files returned
	CheckExistence(files []string) []string

	// Remove removes files from file system
	Remove(files []string)
}

// NewFiler creates new Filer instance
func NewFiler(fs afero.Fs, w io.Writer) Filer {
	return &filer{
		fs: fs,
		w:  w,
	}
}

type filer struct {
	fs afero.Fs
	w  io.Writer
}

// CheckExistence validates files passed to be present in file system
// The list of non exist files returned
func (f *filer) CheckExistence(files []string) []string {
	result := []string{}
	for _, file := range files {
		if _, err := f.fs.Stat(file); os.IsNotExist(err) {
			result = append(result, file)
		}
	}
	return result
}

// Remove removes files from file system
func (f *filer) Remove(files []string) {
	for _, file := range files {
		err := f.fs.Remove(file)
		if err != nil {
			log.Printf("%v\n", err)
		} else {
			_, _ = fmt.Fprintf(f.w, "File: %s removed successfully.\n", file)
		}
	}
}
