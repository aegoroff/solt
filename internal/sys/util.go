package sys

import (
	"bufio"
	"encoding/xml"
	"github.com/spf13/afero"
	"io"
	"log"
	"os"
	"path/filepath"
)

// UnmarshalXMLFrom reads xml from file system and unmarshal it to interface specified
func UnmarshalXMLFrom(path string, fs afero.Fs, result interface{}) error {
	f, err := fs.Open(filepath.Clean(path))
	if err != nil {
		log.Print(err)
		return err
	}
	defer closeResource(f)

	return UnmarshalXML(f, result)
}

// UnmarshalXML reads xml from reader and unmarshal it to interface specified
func UnmarshalXML(r io.Reader, result interface{}) error {
	s := bufio.NewScanner(r)
	var data []byte
	for s.Scan() {
		data = append(data, s.Bytes()...)
	}
	err := xml.Unmarshal(data, result)
	return err
}

// CheckExistence validates files passed to be present in file system
// The list of non exist files returned
func CheckExistence(files []string, fs afero.Fs) []string {
	result := []string{}
	for _, f := range files {
		if _, err := fs.Stat(f); os.IsNotExist(err) {
			result = append(result, f)
		}
	}
	return result
}

func closeResource(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Println(err)
	}
}
