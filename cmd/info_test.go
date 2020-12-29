package cmd

import (
	"bytes"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"solt/cmd/api"
	"solt/solution"
	"testing"
)

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

	w := bytes.NewBufferString("")
	env := api.NewStringEnvironment(w)

	// Act
	_ = Execute(memfs, env, "in", "-p", dir)

	// Assert
	actual := w.String()
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
