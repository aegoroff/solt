package cmd

import (
	"bytes"
	"encoding/xml"
	"fmt"
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/spf13/afero"
	"io"
	"log"
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
		newContent := f.removeRedundantRefsFromProject(project, ends)
		filer.Write(project, newContent)
	}

	f.prn.cprint("Fixed <red>%d</> redundant project references in <red>%d</> projects within solution <red>%s</>\n", invalidRefsCount, len(refs), solution)
}

func (f *sdkProjectsFixer) getElementsEnds(project string, toRemove c9s.StringHashSet) []int64 {
	file, err := f.fs.Open(filepath.Clean(project))
	defer sys.Close(file)
	if err != nil {
		return nil
	}

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
			// If we just read a StartElement token
			if v.Name.Local == "ProjectReference" {
				var prj sdkProjectReference
				// decode a whole chunk of following XML into the variable
				err = decoder.DecodeElement(&prj, &v)
				if err != nil {
					log.Println(err)
				} else {
					referenceFullPath := filepath.Join(pdir, prj.Path)
					if toRemove.Contains(referenceFullPath) {
						ends = append(ends, off)
					}
				}
			}
		}
	}

	return ends
}

func (f *sdkProjectsFixer) removeRedundantRefsFromProject(project string, ends []int64) []byte {
	file, err := f.fs.Open(filepath.Clean(project))
	defer sys.Close(file)
	if err != nil {
		return nil
	}

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
		l := len(portion)
		for i := len(portion) - 2; i >= 0; i-- {
			r, _ := utf8.DecodeRune([]byte{portion[i]})
			if r == '>' || r == '\n' {
				l = i
				break
			}
		}

		portion = portion[:l]
		result = append(result, portion...)
	}

	result = append(result, buf.Next(buf.Len())...)

	return result
}
