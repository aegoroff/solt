package cmd

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NugetCmd_OutputAsSpecified(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	memfs.MkdirAll(dir+"a/Properties", 0755)
	afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent), 0644)
	afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	afero.WriteFile(memfs, dir+"a/packages.config", []byte(packagesConfingContent), 0644)
	afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	appPrinter = newMockPrn()
	appFileSystem = memfs

	// Act
	Execute("nu", "-p", dir, "-r")

	// Assert
	actual := appPrinter.String()
	ass.Equal(`
 a\a
  Package            Version
  -------            -------
  CmdLine            1.0.7.509
  Newtonsoft.Json    12.0.1
`, actual)
}

func Test_NugetCmdFindMismatch_OutputAsSpecified(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	memfs.MkdirAll(dir+"a/Properties", 0755)
	afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent), 0644)
	afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	afero.WriteFile(memfs, dir+"a/packages.config", []byte(packagesConfingContent), 0644)
	afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	appPrinter = newMockPrn()
	appFileSystem = memfs

	// Act
	Execute("nu", "-p", dir, "-m")

	// Assert
	actual := appPrinter.String()
	ass.Equal("", actual)
}

func Test_NugetCmdBySolution_OutputAsSpecified(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	memfs.MkdirAll(dir+"a/Properties", 0755)
	afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent), 0644)
	afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	afero.WriteFile(memfs, dir+"a/packages.config", []byte(packagesConfingContent), 0644)
	afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	appPrinter = newMockPrn()
	appFileSystem = memfs

	// Act
	Execute("nu", "-p", dir)

	// Assert
	actual := appPrinter.String()
	ass.Equal(`
 a\a.sln
  Package            Version
  -------            -------
  CmdLine            1.0.7.509
  Newtonsoft.Json    12.0.1
`, actual)
}
