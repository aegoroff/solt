package msvc

import (
	"bufio"
	"encoding/xml"
	"github.com/aegoroff/dirstat/scan"
	"github.com/spf13/afero"
	"io"
	"log"
	"path/filepath"
)

func unmarshalXMLFrom(path string, fs afero.Fs, result interface{}) error {
	f, err := fs.Open(filepath.Clean(path))
	if err != nil {
		log.Print(err)
		return err
	}
	defer scan.Close(f)
	s, err := f.Stat()
	if err != nil {
		log.Print(err)
		return err
	}

	return unmarshalXML(f, result, int(s.Size()))
}

func unmarshalXML(r io.Reader, result interface{}, sz int) error {
	s := bufio.NewScanner(r)
	data := make([]byte, sz)
	start := 0
	for s.Scan() {
		start += copy(data[start:], s.Bytes())
	}
	err := xml.Unmarshal(data, result)
	return err
}
