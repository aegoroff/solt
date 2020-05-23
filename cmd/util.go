package cmd

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"github.com/spf13/afero"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type filesystemItem struct {
	dir   string
	entry os.FileInfo
}

func outputSortedMap(writer io.Writer, itemsMap map[string][]string, keyPrefix string) {
	var keys []string
	for k := range itemsMap {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(writer, "\n%s: %s\n", keyPrefix, k)
		sortAndOutput(writer, itemsMap[k])
	}
}

func sortAndOutput(writer io.Writer, items []string) {
	sort.Strings(items)
	for _, item := range items {
		fmt.Fprintf(writer, " %s\n", item)
	}
}

func unmarshalXmlFrom(path string, fs afero.Fs, result interface{}) error {
	f, err := fs.Open(filepath.Clean(path))
	if err != nil {
		log.Print(err)
		return err
	}
	defer closeResource(f)

	return unmarshalXml(f, result)
}

func unmarshalXml(r io.Reader, result interface{}) error {
	s := bufio.NewScanner(r)
	var data []byte
	for s.Scan() {
		data = append(data, s.Bytes()...)
	}
	err := xml.Unmarshal(data, result)
	return err
}

func walkDirBreadthFirst(path string, fs afero.Fs, results chan<- filesystemItem) {
	defer close(results)

	var wg sync.WaitGroup
	var mu sync.RWMutex
	queue := make([]string, 0)

	queue = append(queue, path)

	ql := len(queue)

	for ql > 0 {
		// Peek
		mu.RLock()
		currentDir := queue[0]
		mu.RUnlock()

		wg.Add(1)
		go func(d string) {
			defer wg.Done()

			entries := dirents(d, fs)

			if entries == nil {
				return
			}

			for _, entry := range entries {
				if entry.IsDir() {
					// Queue subdirs to walk in a queue
					subdir := filepath.Join(d, entry.Name())

					// Push
					mu.Lock()
					queue = append(queue, subdir)
					mu.Unlock()
				} else {
					// Send file to channel
					results <- filesystemItem{
						dir:   d,
						entry: entry,
					}
				}
			}
		}(currentDir)

		// Pop
		mu.Lock()
		queue = queue[1:]
		ql = len(queue)
		mu.Unlock()

		if ql == 0 {
			// Waiting pending goroutines
			wg.Wait()

			mu.RLock()
			ql = len(queue)
			mu.RUnlock()
		}
	}
}

func dirents(path string, fs afero.Fs) []os.FileInfo {
	entries, err := ReadDir(path, fs)
	if err != nil {
		return nil
	}

	return entries
}

// ReadDir reads the directory named by dirname and returns
// a list of directory entries sorted by filename.
func ReadDir(dirname string, fs afero.Fs) ([]os.FileInfo, error) {
	f, err := fs.Open(dirname)
	if err != nil {
		return nil, err
	}

	defer closeResource(f)

	list, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, nil
}

func closeResource(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Println(err)
	}
}
