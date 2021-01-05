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
	switch v := t.(type) {
	case xml.StartElement:
		if v.Name.Local == "ProjectReference" {
			var prj projectReference
			// decode a whole chunk of following XML into the variable
			offBefore := d.InputOffset()
			_ = d.DecodeElement(&prj, &v)
			offAfter := d.InputOffset()

			r.addEnds(prj, offBefore, offAfter)
		}
	}
}

func (r *elementRemover) addEnds(prj projectReference, offBefore int64, offAfter int64) {
	fullPath := filepath.Join(r.dir, prj.path())
	if r.toRemove.Contains(fullPath) {
		r.ends = append(r.ends, offBefore)
		if offAfter > offBefore {
			r.ends = append(r.ends, offAfter)
		}
	}
}
