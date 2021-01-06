package sys

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/spf13/afero"
	"io"
)

// XMLDecoder defines XML decoder that works as SAX parser
type XMLDecoder struct {
	w io.Writer
}

// NewXMLDecoder creates new decoder instance
func NewXMLDecoder(w io.Writer) *XMLDecoder {
	return &XMLDecoder{w: w}
}

// DecodeFn provides decoder token function prototype
type DecodeFn func(d *xml.Decoder, t xml.Token)

// Decode decodes xml file execute decoders on each token
func (x *XMLDecoder) Decode(rdr io.Reader, decoders ...DecodeFn) {
	decoder := xml.NewDecoder(rdr)
	for {
		t, err := decoder.Token()
		if t == nil {
			x.print(err)
			break
		}

		for _, fn := range decoders {
			fn(decoder, t)
		}
	}
}

// UnmarshalFrom unmarshal whole xml file using path specified
func (x *XMLDecoder) UnmarshalFrom(path string, fs afero.Fs, result interface{}) error {
	filer := NewFiler(fs, x.w)
	b, err := filer.Read(path)
	if err != nil {
		return err
	}

	r := bytes.NewReader(b)
	return x.Unmarshal(r, result)
}

// Unmarshal unmarshal whole xml file using reader specified
func (*XMLDecoder) Unmarshal(r io.Reader, result interface{}) error {
	return xml.NewDecoder(r).Decode(result)
}

func (x *XMLDecoder) print(err error) {
	if err != nil && err != io.EOF && x.w != nil {
		_, _ = fmt.Fprintln(x.w, err)
	}
}
