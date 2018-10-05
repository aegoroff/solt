package solution

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
)

const visualStudioVersionKey = "VisualStudioVersion"
const minimumVisualStudioVersionKey = "MinimumVisualStudioVersion"

var (
    visualStudioVersion        string
    minimumVisualStudioVersion string
    comment                    string
    words                      []string
)

var projects []*Project
var globalSections []*Section
var currentSectionType string

func (l *lexer) Lex(lval *yySymType) int {
    v := l.nextItem()
    if v.tok == itemEOF {
        lval.tok = 0
    } else {
        lval.tok = v.tok
    }
    lval.str = v.str
    lval.line = v.line
    lval.yys = v.yys
    //fmt.Printf("%s:%q\n",lval.tok, lval.str)
    return int(lval.tok)
}

func (l *lexer) Error(e string) {
    // TODO: fix this spike
    if !l.atEOF {
        log.Print(e)
    } else if yyDebug >= 1 {
        log.Print(e)
    }
}

// Parses visual studio solution file specified by path
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
        br.UnreadRune() // Not a BOM -- put the rune back
    }

    bs := bufio.NewScanner(br)
    bs.Split(bufio.ScanRunes)
    sb := strings.Builder{}

    for bs.Scan() {
        sb.WriteString(bs.Text())
    }

    str := sb.String()

    parse(str)

    sol := Solution{
        GlobalSections:             globalSections,
        Projects:                   projects,
        MinimumVisualStudioVersion: minimumVisualStudioVersion,
        VisualStudioVersion:        visualStudioVersion,
        Comment:                    comment,
        Header:                     strings.Join(words, " "),
    }

    return &sol, nil
}

func parse(str string) {
    //yyDebug = 3
    projects = []*Project{}
    globalSections = []*Section{}
    minimumVisualStudioVersion = ""
    visualStudioVersion = ""
    comment = ""
    words = []string{}
    yyErrorVerbose = true
    lx := newLexer(str)
    yyParse(lx)
}

func onProject(projectType, name, path, id string) {
    p := Project{
        TypeId: strings.Trim(projectType, " "),
        Name:   strings.Trim(name, " "),
        Path:   strings.Trim(path, " "),
        Id:     strings.Trim(id, " "),
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
        words = append(words, fmt.Sprintf("%s,", value))
    } else {
        words = append(words, value)
    }
}

const projectSection = "ProjectSection"

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
