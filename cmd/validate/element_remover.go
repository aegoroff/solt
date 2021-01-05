package validate

import (
	"encoding/xml"
	c9s "github.com/aegoroff/godatastruct/collections"
	"path/filepath"
)

type elementRemover struct {
	toRemove c9s.StringHashSet
	ends     []int64
	dir      string
}

func newElementRemover(project string, toRemove c9s.StringHashSet) *elementRemover {
	return &elementRemover{
		toRemove: toRemove,
		ends:     make([]int64, 0),
		dir:      filepath.Dir(project),
	}
}

func (r *elementRemover) decode(d *xml.Decoder, t xml.Token) {
	off := d.InputOffset()

	switch v := t.(type) {
	case xml.StartElement:
		if v.Name.Local == "ProjectReference" {
			var prj projectReference
			// decode a whole chunk of following XML into the variable
			_ = d.DecodeElement(&prj, &v)
			offAfter := d.InputOffset()

			referenceFullPath := filepath.Join(r.dir, prj.path())
			if r.toRemove.Contains(referenceFullPath) {
				r.ends = append(r.ends, off)
				if offAfter > off {
					r.ends = append(r.ends, offAfter)
				}
			}
		}
	}
}
