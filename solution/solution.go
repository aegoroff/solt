package solution

import (
	"bufio"
	"io"
	"log"
	"strings"
)

const (
	visualStudioVersionKey        = "VisualStudioVersion"
	minimumVisualStudioVersionKey = "MinimumVisualStudioVersion"
	projectSection                = "ProjectSection"
)

var (
	visualStudioVersion        string
	minimumVisualStudioVersion string
	comment                    string
	words                      []string
	projects                   []*Project
	globalSections             []*Section
	currentSectionType         string
)

func (lx *lexer) Lex(lval *yySymType) int {
	v := lx.nextItem()
	if v.tok == itemEOF {
		return 0
	}
	if lx.debug {
		log.Println(v.tok.String())
	}
	lval.tok = v.tok
	lval.str = v.str
	lval.line = v.line
	lval.yys = v.yys
	return int(lval.tok)
}

func (*lexer) Error(e string) {
	log.Print(e)
}

// Parse parses visual studio solution file specified by io.Reader
func Parse(rdr io.Reader) (*Solution, error) {
	br := bufio.NewReader(rdr)
	r, _, err := br.ReadRune()
	if err != nil {
		return nil, err
	}
	if r != '\uFEFF' {
		err = br.UnreadRune() // Not a BOM -- put the rune back
		if err != nil {
			return nil, err
		}
	}

	bs := bufio.NewScanner(br)
	bs.Split(bufio.ScanRunes)
	sb := strings.Builder{}

	for bs.Scan() {
		_, err = sb.WriteString(bs.Text())
		if err != nil {
			return nil, err
		}
	}

	str := sb.String()

	sol := parse(str, false)

	return sol, nil
}

func parse(str string, debug bool) *Solution {
	projects = []*Project{}
	globalSections = []*Section{}
	minimumVisualStudioVersion = ""
	visualStudioVersion = ""
	comment = ""
	words = []string{}
	yyErrorVerbose = true
	lx := newLexer(str, debug)
	yyParse(lx)

	return &Solution{
		GlobalSections:             globalSections,
		Projects:                   projects,
		MinimumVisualStudioVersion: minimumVisualStudioVersion,
		VisualStudioVersion:        visualStudioVersion,
		Comment:                    comment,
		Header:                     strings.Join(words, " "),
	}
}

func onProject(projectType, name, path, id string) {
	p := Project{
		TypeID: projectType,
		Name:   name,
		Path:   path,
		ID:     id,
	}
	if v, ok := ProjectsGuids[p.TypeID]; ok {
		p.Type = v
	}

	projects = append(projects, &p)
}

func onVersion(key, value string) {
	switch key {
	case minimumVisualStudioVersionKey:
		minimumVisualStudioVersion = value
	case visualStudioVersionKey:
		visualStudioVersion = value
	}
}

func onComment(value string) {
	comment = value
}

func onWord(value string) {
	if value == "File" {
		// HACK
		words = append(words, value+",")
	} else {
		words = append(words, value)
	}
}

func onSection(sectionType, name, stage string) {
	s := Section{Name: name, Stage: stage}
	currentSectionType = sectionType
	if sectionType == projectSection {
		projects[len(projects)-1].Sections = append(projects[len(projects)-1].Sections, &s)
	} else {
		globalSections = append(globalSections, &s)
	}
}

func onSectionItem(key, value string) {
	si := SectionItem{Key: key, Value: value}
	if currentSectionType == projectSection {
		prj := projects[len(projects)-1]
		sections := prj.Sections
		sectionIx := len(sections) - 1

		prj.Sections[sectionIx].Items = append(prj.Sections[sectionIx].Items, &si)
	} else {
		globalSections[len(globalSections)-1].Items = append(globalSections[len(globalSections)-1].Items, &si)
	}
}
