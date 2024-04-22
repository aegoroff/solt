package va

import (
	"encoding/xml"
	"path/filepath"
	"solt/internal/sys"

	c9s "github.com/aegoroff/godatastruct/collections"
)

type projectReference struct {
	Path string `xml:"Include,attr"`
}

func (r *projectReference) path() string {
	return sys.ToValidPath(r.Path)
}

type elementEndDetector struct {
	filter c9s.HashSet[string]
	ends   []int64
	dir    string
}

func newElementEndDetector(project string, filter c9s.HashSet[string]) *elementEndDetector {
	return &elementEndDetector{
		filter: filter,
		ends:   make([]int64, 0),
		dir:    filepath.Dir(project),
	}
}

func (r *elementEndDetector) decode(d *xml.Decoder, t xml.Token) {
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

func (r *elementEndDetector) addEnds(prj projectReference, offBefore int64, offAfter int64) {
	fullPath := filepath.Join(r.dir, prj.path())
	if r.filter.Contains(fullPath) {
		r.ends = append(r.ends, offBefore)
		if offAfter > offBefore {
			r.ends = append(r.ends, offAfter)
		}
	}
}
