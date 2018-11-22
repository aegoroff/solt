package cmd

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func outputSortedMapToStdout(itemsMap map[string][]string, keyPrefix string) {
	outputSortedMap(os.Stdout, itemsMap, keyPrefix)
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

func sortAndOutputToStdout(items []string) {
	sortAndOutput(os.Stdout, items)
}

func sortAndOutput(writer io.Writer, items []string) {
	sort.Strings(items)
	for _, item := range items {
		fmt.Fprintf(writer, " %s\n", item)
	}
}

func unmarshalXmlFrom(path string, result interface{}) error {
	f, err := os.Open(path)
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

func walkDirBreadthFirst(path string, action func(parent string, entry os.FileInfo)) {
	queue := make([]string, 0)

	queue = append(queue, path)

	for len(queue) > 0 {
		curr := queue[0]

		for _, entry := range dirents(curr) {
			action(curr, entry)
			if entry.IsDir() {
				queue = append(queue, filepath.Join(curr, entry.Name()))
			}
		}

		queue = queue[1:]
	}
}

func dirents(path string) []os.FileInfo {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return nil
	}

	return entries
}

func closeResource(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Println(err)
	}
}
