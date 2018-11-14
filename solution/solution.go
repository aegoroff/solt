package solution

import (
	"bufio"
	"log"
	"os"
	"strings"
)

const visualStudioVersionKey = "VisualStudioVersion"
const minimumVisualStudioVersionKey = "MinimumVisualStudioVersion"
const projectSection = "ProjectSection"

var (
	visualStudioVersion        string
	minimumVisualStudioVersion string
	comment                    string
	words                      []string
	projects                   []*Project
	globalSections             []*Section
	currentSectionType         string
)

func (l *lexer) Lex(lval *yySymType) int {
	v := l.nextItem()
	if v.tok == itemEOF {
		return 0
	}
	lval.tok = v.tok
	lval.str = v.str
	lval.line = v.line
	lval.yys = v.yys
	return int(lval.tok)
}

func (l *lexer) Error(e string) {
	log.Print(e)
}

// Parse parses visual studio solution file specified by path
func Parse(solutionPath string) (*Solution, error) {

	f, err := os.Open(solutionPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	br := bufio.NewReader(f)
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
		sb.WriteString(bs.Text())
	}

	str := sb.String()

	sol := parse(str)

	return sol, nil
}

func parse(str string) *Solution {
	projects = []*Project{}
	globalSections = []*Section{}
	minimumVisualStudioVersion = ""
	visualStudioVersion = ""
	comment = ""
	words = []string{}
	yyErrorVerbose = true
	lx := newLexer(str)
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
		TypeId: projectType,
		Name:   name,
		Path:   path,
		Id:     id,
	}
	if v, ok := ProjectsGuids[p.TypeId]; ok {
		p.Type = v
	}

	projects = append(projects, &p)
}

func onVersion(key, value string) {
	switch key {
	case minimumVisualStudioVersionKey:
		minimumVisualStudioVersion = value
		break
	case visualStudioVersionKey:
		visualStudioVersion = value
		break
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
