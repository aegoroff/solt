package validate

import (
	"bytes"
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/spf13/afero"
	"solt/cmd/api"
	"solt/internal/sys"
	"unicode/utf8"
)

type fixer struct {
	prn   api.Printer
	fs    afero.Fs
	filer sys.Filer
	w     api.Writable
}

func newFixer(p api.Printer, w api.Writable, fs afero.Fs) actioner {
	return &fixer{
		prn:   p,
		fs:    fs,
		w:     w,
		filer: sys.NewFiler(fs, w.Writer()),
	}
}

func (f *fixer) action(path string, refs map[string]c9s.StringHashSet) {
	if len(refs) == 0 {
		return
	}

	invalidRefsCount := 0
	for project, rrs := range refs {
		invalidRefsCount += rrs.Count()
		buf, err := f.filer.Read(project)

		if err == nil {
			ends := f.getElementsEnds(buf, project, rrs)
			newContent := f.getNewFileContent(buf, ends)
			f.filer.Write(project, newContent)
		}
	}

	const mf = "Fixed <red>%d</> redundant project references in <red>%d</> projects within solution <red>%s</>\n"
	f.prn.Cprint(mf, invalidRefsCount, len(refs), path)
}

func (f *fixer) getElementsEnds(buf *bytes.Buffer, project string, toRemove c9s.StringHashSet) []int64 {
	ed := newElementEndDetector(project, toRemove)

	decoder := api.NewXMLDecoder(f.w.Writer())
	r := bytes.NewReader(buf.Bytes())
	decoder.Decode(r, ed.decode)

	return ed.ends
}

func (f *fixer) getNewFileContent(buf *bytes.Buffer, ends []int64) []byte {
	result := make([]byte, 0, buf.Len())
	start := 0
	for _, end := range ends {
		n := int(end) - start
		start = int(end)

		portion := buf.Next(n)
		l := fallback(portion)

		portion = portion[:l]
		result = append(result, portion...)
	}

	result = append(result, buf.Next(buf.Len())...)

	return result
}

func fallback(data []byte) int {
	for i := len(data) - 2; i >= 0; i-- {
		l, ok := stopFallback(data, i)
		if ok {
			return l
		}
	}
	return 0
}

func stopFallback(data []byte, current int) (int, bool) {
	r := getRune(data, current)

	ix := len(data)
	ok := false

	switch r {
	case '>':
		// Keep > so as not to break xml syntax
		ix = current + 1
		ok = true
		break
	case '\n':
		ix = current
		ok = true
		if current > 0 && getRune(data, current-1) == '\r' {
			// windows case (\r\n as line ending) - remove \r to keep correctness
			ix = current - 1
		}
		break
	}

	return ix, ok
}

func getRune(data []byte, current int) rune {
	r, _ := utf8.DecodeRune([]byte{data[current]})
	return r
}
