package fw

import (
	"io"
	"solt/internal/out"
	"solt/internal/sys"
	"solt/internal/ux"

	"github.com/aegoroff/godatastruct/collections"
	"github.com/spf13/afero"
)

type exister struct {
	filer        Filer
	missingFiles map[string][]string
	missingPaths collections.HashSet[string]
}

// NewExister creates new Exister instance
func NewExister(fs afero.Fs, w io.Writer) Exister {
	return &exister{
		missingFiles: make(map[string][]string),
		missingPaths: make(collections.HashSet[string]),
		filer:        sys.NewFiler(fs, w),
	}
}

// Validate validates whether files from container exist in filesystem
func (e *exister) Validate(root string, paths []string) {
	missingFiles := e.filer.CheckExistence(paths)
	e.missingPaths.AddRange(missingFiles...)

	if len(missingFiles) > 0 {
		e.missingFiles[root] = append(e.missingFiles[root], missingFiles...)
	}
}

// MissingCount gets the number of missing items
func (e *exister) MissingCount() int64 {
	return int64(e.missingPaths.Count())
}

// Print outputs missing files in
func (e *exister) Print(p out.Printer, title string, container string) {
	if len(e.missingFiles) > 0 {
		p.Println()
		p.Cprint(title)
		p.Println()
	}

	s := ux.NewScreener(p)
	s.WriteMap(e.missingFiles, container)
}
