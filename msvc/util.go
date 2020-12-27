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

	return unmarshalXML(f, result)
}

func unmarshalXML(r io.Reader, result interface{}) error {
	s := bufio.NewScanner(r)
	var data []byte
	for s.Scan() {
		data = append(data, s.Bytes()...)
	}
	err := xml.Unmarshal(data, result)
	return err
}
