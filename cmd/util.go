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
)

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

func walkDirBreadthFirst(path string, fs afero.Fs, action func(parent string, entry os.FileInfo)) {
	queue := make([]string, 0)

	queue = append(queue, path)

	for len(queue) > 0 {
		curr := queue[0]

		for _, entry := range dirents(curr, fs) {
			action(curr, entry)
			if entry.IsDir() {
				queue = append(queue, filepath.Join(curr, entry.Name()))
			}
		}

		queue = queue[1:]
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
