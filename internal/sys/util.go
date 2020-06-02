package sys

import (
	"bufio"
	"encoding/xml"
	"github.com/spf13/afero"
	"io"
	"log"
	"path/filepath"
)

// UnmarshalXmlFrom reads xml from file system and unmarshal it to interface specified
func UnmarshalXmlFrom(path string, fs afero.Fs, result interface{}) error {
	f, err := fs.Open(filepath.Clean(path))
	if err != nil {
		log.Print(err)
		return err
	}
	defer closeResource(f)

	return UnmarshalXml(f, result)
}

// UnmarshalXml reads xml from reader and unmarshal it to interface specified
func UnmarshalXml(r io.Reader, result interface{}) error {
	s := bufio.NewScanner(r)
	var data []byte
	for s.Scan() {
		data = append(data, s.Bytes()...)
	}
	err := xml.Unmarshal(data, result)
	return err
}

func closeResource(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Println(err)
	}
}
