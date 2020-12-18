package sys

import (
	"fmt"
	"github.com/spf13/afero"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Filer defines module that works with files
type Filer interface {
	// CheckExistence validates files passed to be present in file system
	// The list of non exist files returned
	CheckExistence(files []string) []string

	// Remove removes files from file system
	Remove(files []string)

	// Write writes new file content
	Write(path string, bytes []byte)
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
	result := make([]string, 0)
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

func (f *filer) Write(path string, bytes []byte) {
	fi, err := f.fs.Create(filepath.Clean(path))
	defer Close(fi)
	if err != nil || bytes == nil {
		return
	}
	_, err = fi.Write(bytes)
	if err != nil {
		_, _ = fmt.Fprintf(f.w, "%v\n", err)
	}
}
