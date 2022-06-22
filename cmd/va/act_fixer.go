package va

import (
	"bytes"
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/spf13/afero"
	"solt/internal/fw"
	"solt/internal/out"
	"solt/internal/sys"
	"unicode/utf8"
)

type fixer struct {
	prn   out.Printer
	fs    afero.Fs
	filer fw.Filer
	w     out.Writable
}

func newFixer(p out.Printer, w out.Writable, fs afero.Fs) actioner {
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
		data, err := f.filer.Read(project)

		if err == nil {
			ends := f.getElementsEnds(data, project, rrs)
			newContent := getNewFileContent(data, ends)
			f.filer.Write(project, newContent)
		}
	}

	const mf = "Fixed <red>%d</> redundant project references in <red>%d</> projects within solution <red>%s</>\n"
	f.prn.Cprint(mf, invalidRefsCount, len(refs), path)
}

func (f *fixer) getElementsEnds(data []byte, project string, toRemove c9s.StringHashSet) []int64 {
	ed := newElementEndDetector(project, toRemove)

	decoder := sys.NewXMLDecoder(f.w.Writer())
	r := bytes.NewReader(data)
	decoder.Decode(r, ed.decode)

	return ed.ends
}

func getNewFileContent(data []byte, ends []int64) []byte {
	result := make([]byte, 0, len(data))
	start := 0
	buf := bytes.NewBuffer(data)
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
	case '\n':
		ix = current
		ok = true
		if current > 0 && getRune(data, current-1) == '\r' {
			// windows case (\r\n as line ending) - remove \r to keep correctness
			ix = current - 1
		}
	}

	return ix, ok
}

func getRune(data []byte, current int) rune {
	r, _ := utf8.DecodeRune([]byte{data[current]})
	return r
}
