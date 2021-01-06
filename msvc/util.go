package msvc

import (
	"bytes"
	"encoding/xml"
	"github.com/spf13/afero"
	"io"
	"os"
	"solt/internal/sys"
)

func unmarshalXMLFrom(path string, fs afero.Fs, result interface{}) error {
	filer := sys.NewFiler(fs, os.Stderr)
	b, err := filer.Read(path)
	if err != nil {
		return err
	}

	r := bytes.NewReader(b)
	return unmarshalXML(r, result)
}

func unmarshalXML(r io.Reader, result interface{}) error {
	return xml.NewDecoder(r).Decode(result)
}
