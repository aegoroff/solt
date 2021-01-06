package sys

import (
	"bytes"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestXMLDecoder_Decode(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	const invalid = "<Project><Item></Project>"
	buf := bytes.NewBufferString(invalid)
	out := bytes.NewBuffer(nil)
	d := NewXMLDecoder(out)
	i := 0

	// Act
	d.Decode(buf, func(d *xml.Decoder, t xml.Token) {
		i++
	})

	// Assert
	ass.True(i > 0)
	ass.Equal("XML syntax error on line 1: element <Item> closed by </Project>\n", out.String())
}
