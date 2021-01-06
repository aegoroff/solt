package sys

import (
	"bytes"
	"fmt"
	"github.com/aegoroff/dirstat/scan"
	"github.com/spf13/afero"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// Filer defines module that works with files
type Filer interface {
	// CheckExistence validates files passed to be present in file system
	// The list of non exist files returned
	CheckExistence(files []string) []string

	// Remove removes files from file system
	Remove(files []string)

	// Write writes new file content
	Write(path string, content []byte)

	// Read reads file content
	Read(path string) ([]byte, error)
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
	var mu sync.RWMutex
	result := make([]string, 0)
	var restrict = make(chan struct{}, 32)
	defer close(restrict)

	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, file := range files {
		restrict <- struct{}{}
		go func(file string, restrict chan struct{}) {
			defer wg.Done()
			defer func() { <-restrict }()

			if f.fileNotExists(file) {
				mu.Lock()
				result = append(result, file)
				mu.Unlock()
			}
		}(file, restrict)
	}

	wg.Wait()
	return result
}

func (f *filer) fileNotExists(path string) bool {
	_, err := f.fs.Stat(path)
	return os.IsNotExist(err)
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

func (f *filer) Write(path string, content []byte) {
	if content == nil {
		return
	}

	fi, err := f.fs.Create(filepath.Clean(path))
	if err != nil {
		return
	}
	defer scan.Close(fi)

	_, err = fi.Write(content)
	if err != nil {
		_, _ = fmt.Fprintf(f.w, "%v\n", err)
	}
}

func (f *filer) Read(path string) ([]byte, error) {
	file, err := f.fs.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	defer scan.Close(file)

	s, err := file.Stat()
	if err != nil {
		return nil, err
	}

	b := make([]byte, 0, s.Size())

	buf := bytes.NewBuffer(b)
	_, err = io.Copy(buf, file)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
