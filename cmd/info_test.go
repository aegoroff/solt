package cmd

import (
	"bytes"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_InfoCmd_InfoAsSpecified(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	memfs.MkdirAll(dir+"a/Properties", 0755)
	afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent), 0644)
	afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	buf := bytes.NewBufferString("")

	appWriter = buf
	appFileSystem = memfs

	// Act
	rootCmd.SetArgs([]string{"in", "-p", dir})
	rootCmd.Execute()

	// Assert
	actual := buf.String()
	ass.Equal(`  Header                            Microsoft Visual Studio Solution File, Format Version 12.00
  Product                           # Visual Studio Version 16
  Visial Studion Version            16.0.30104.148
  Minimum Visial Studion Version    10.0.40219.1
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
`, actual)
}
