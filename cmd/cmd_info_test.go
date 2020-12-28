package cmd

import (
	"bytes"
	"fmt"
	"github.com/gookit/color"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"io"
	"solt/solution"
	"testing"
	"text/tabwriter"
)

type mockprn struct {
	tw *tabwriter.Writer
	w  *bytes.Buffer
}

func (m *mockprn) String() string {
	return m.w.String()
}

func newMockPrn() *mockprn {
	w := bytes.NewBufferString("")
	tw := new(tabwriter.Writer).Init(w, 0, 8, 4, ' ', 0)

	p := mockprn{
		tw: tw,
		w:  w,
	}
	return &p
}

func (m *mockprn) Tprint(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(m.tw, format, a...)
}

func (m *mockprn) Cprint(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	_, _ = fmt.Fprintf(m.w, str)
}

func (m *mockprn) Writer() io.Writer { return m.w }

func (m *mockprn) Twriter() *tabwriter.Writer { return m.tw }

func (*mockprn) SetColor(_ color.Color) {}

func (*mockprn) ResetColor() {}

func (m *mockprn) Flush() {
	_ = m.tw.Flush()
}

func Test_InfoCmd_InfoAsSpecified(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	_ = memfs.MkdirAll(dir+"a/Properties", 0755)
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/packages.config", []byte(packagesConfingContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	p := newMockPrn()

	// Act
	_ = execute(memfs, p, "in", "-p", dir)

	// Assert
	actual := p.w.String()
	ass.Equal(solution.ToValidPath(` a\a.sln
  Header                           Microsoft Visual Studio Solution File, Format Version 12.00
  Product                          # Visual Studio Version 16
  Visual Studio Version            16.0.30104.148
  Minimum Visual Studio Version    10.0.40219.1
  Project type    Count
  ------------    -----
  C#              1
  Configuration
  ------------
  Debug
  Release
  Platform
  --------
  Any CPU

`), actual)
}
