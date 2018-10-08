package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"github.com/dustin/go-humanize"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
)

// printMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("\nAlloc = %s", humanize.IBytes(m.Alloc))
	fmt.Printf("\tTotalAlloc = %s", humanize.IBytes(m.TotalAlloc))
	fmt.Printf("\tSys = %s", humanize.IBytes(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func sortAndOutput(items []string) {
	sort.Strings(items)
	for _, item := range items {
		fmt.Printf(" %s\n", item)
	}
}

func unmarshalXml(path string, result interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	var data []byte
	for s.Scan() {
		data = append(data, s.Bytes()...)
	}
	err = xml.Unmarshal(data, result)
	if err != nil {
		return err
	}

	return nil
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
