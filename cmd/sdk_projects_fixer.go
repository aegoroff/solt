package cmd

import (
	"bytes"
	"encoding/xml"
	"fmt"
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/spf13/afero"
	"io"
	"path/filepath"
	"solt/internal/sys"
	"unicode/utf8"
)

type sdkProjectsFixer struct {
	prn printer
	fs  afero.Fs
}

func newsdkProjectsFixer(p printer, fs afero.Fs) sdkModuleHandler {
	return &sdkProjectsFixer{
		prn: p,
		fs:  fs,
	}
}

func (f *sdkProjectsFixer) handle(solution string, refs map[string]c9s.StringHashSet) {
	if len(refs) == 0 {
		return
	}

	filer := sys.NewFiler(f.fs, f.prn.writer())

	invalidRefsCount := 0
	for project, rrs := range refs {
		invalidRefsCount += rrs.Count()
		ends := f.getElementsEnds(project, rrs)
		newContent := f.getNewFileContent(project, ends)
		filer.Write(project, newContent)
	}

	const mf = "Fixed <red>%d</> redundant project references in <red>%d</> projects within solution <red>%s</>\n"
	f.prn.cprint(mf, invalidRefsCount, len(refs), solution)
}

func (f *sdkProjectsFixer) getElementsEnds(project string, toRemove c9s.StringHashSet) []int64 {
	file, err := f.fs.Open(filepath.Clean(project))
	if err != nil {
		return nil
	}
	defer sys.Close(file)

	decoder := xml.NewDecoder(file)
	pdir := filepath.Dir(project)

	ends := make([]int64, 0)
	for {
		t, err := decoder.Token()
		if t == nil {
			if err != nil && err != io.EOF {
				fmt.Println(err)
			}
			break
		}

		off := decoder.InputOffset()

		switch v := t.(type) {
		case xml.StartElement:
			if v.Name.Local == "ProjectReference" {
				var prj sdkProjectReference
				// decode a whole chunk of following XML into the variable
				_ = decoder.DecodeElement(&prj, &v)
				offAfter := decoder.InputOffset()
				referenceFullPath := filepath.Join(pdir, prj.Path)
				if toRemove.Contains(referenceFullPath) {
					ends = append(ends, off)
					if offAfter > off {
						ends = append(ends, offAfter)
					}
				}
			}
		}
	}

	return ends
}

func (f *sdkProjectsFixer) getNewFileContent(project string, ends []int64) []byte {
	file, err := f.fs.Open(filepath.Clean(project))
	if err != nil {
		return nil
	}
	defer sys.Close(file)

	buf := bytes.NewBuffer(nil)
	written, err := io.Copy(buf, file)

	if err != nil {
		return nil
	}

	result := make([]byte, 0, written)
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
	stop := r == '>' || r == '\n'
	l := -1
	if stop {
		l = current
	}

	// remove \n so as not to have empty line in file
	if r == '\n' && current > 0 {
		prev := getRune(data, current-1)
		if prev == '\r' || prev == '\n' {
			l = current - 1
		}
	}

	if r == '>' {
		l = current + 1
	}

	return l, stop
}

func getRune(data []byte, current int) rune {
	r, _ := utf8.DecodeRune([]byte{data[current]})
	return r
}
