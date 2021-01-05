package api

import (
	"encoding/xml"
	"fmt"
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
			if err != nil && err != io.EOF {
				_, _ = fmt.Fprintln(x.w, err)
			}
			break
		}

		for _, fn := range decoders {
			fn(decoder, t)
		}
	}
}
