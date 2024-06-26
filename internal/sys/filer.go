package sys

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/aegoroff/dirstat/scan"
	"github.com/aegoroff/godatastruct/collections"
	"github.com/spf13/afero"
)

// NewFiler creates new Filer instance
func NewFiler(fs afero.Fs, w io.Writer) *Filer {
	return &Filer{
		fs: fs,
		w:  w,
	}
}

// Filer provides methods to work with files i.e. reading, writing
// checking existence and remove
type Filer struct {
	fs afero.Fs
	w  io.Writer
}

// CheckExistence validates files passed to be present in file system
// The list of non exist files returned
func (f *Filer) CheckExistence(files []string) []string {
	var mu sync.RWMutex
	result := make([]string, 0)
	var restrict = make(chan struct{}, 32)
	defer close(restrict)
	notExitsDirs := collections.NewHashSet[string]()

	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, file := range files {
		restrict <- struct{}{}
		go func(file string, restrict chan struct{}) {
			defer wg.Done()
			defer func() { <-restrict }()
			dir := filepath.Dir(file)
			if notExitsDirs.Contains(dir) {
				mu.Lock()
				result = append(result, file)
				mu.Unlock()
				return
			}

			if f.fileNotExists(file) {
				if f.fileNotExists(dir) {
					mu.Lock()
					notExitsDirs.Add(dir)
					mu.Unlock()
				}
				mu.Lock()
				result = append(result, file)
				mu.Unlock()
			}
		}(file, restrict)
	}

	wg.Wait()
	return result
}

func (f *Filer) fileNotExists(path string) bool {
	_, err := f.fs.Stat(path)
	return os.IsNotExist(err)
}

// Remove removes files from file system
func (f *Filer) Remove(files []string) {
	for _, file := range files {
		err := f.fs.Remove(file)
		if err != nil {
			log.Printf("%v\n", err)
		} else {
			_, _ = fmt.Fprintf(f.w, "File: %s removed successfully.\n", file)
		}
	}
}

func (f *Filer) Write(path string, content []byte) {
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

func (f *Filer) Read(path string) ([]byte, error) {
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
