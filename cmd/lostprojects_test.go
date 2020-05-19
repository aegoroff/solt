package cmd

import (
	"bytes"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_FindLostProjectsCmd_NoLostProjectsFound(t *testing.T) {
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
	rootCmd.SetArgs([]string{"lp", "-p", dir})
	rootCmd.Execute()

	// Assert
	actual := buf.String()
	ass.Equal(``, actual)
}

func Test_FindLostProjectsCmdLostProjectsInTheSameDir_LostProjectsFound(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	memfs.MkdirAll(dir+"a/Properties", 0755)
	afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent), 0644)
	afero.WriteFile(memfs, dir+"a/a1.csproj", []byte(testProjectContent2), 0644)
	afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	buf := bytes.NewBufferString("")

	appWriter = buf
	appFileSystem = memfs

	// Act
	rootCmd.SetArgs([]string{"lp", "-p", dir})
	rootCmd.Execute()

	// Assert
	actual := buf.String()
	ass.Equal(``, actual)
}

func Test_FindLostProjectsCmdLostProjectsInOtherDir_LostProjectsFound(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	memfs.MkdirAll(dir+"a/Properties", 0755)
	afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent), 0644)
	afero.WriteFile(memfs, dir+"a1/a1.csproj", []byte(testProjectContent2), 0644)
	afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	buf := bytes.NewBufferString("")

	appWriter = buf
	appFileSystem = memfs

	// Act
	rootCmd.SetArgs([]string{"lp", "-p", dir})
	rootCmd.Execute()

	// Assert
	actual := buf.String()
	ass.Equal(` a\a1\a1.csproj
`, actual)
}

func Test_FindLostProjectsCmdUnexistProjects_LostProjectsFound(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	memfs.MkdirAll(dir+"a/Properties", 0755)
	afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	buf := bytes.NewBufferString("")

	appWriter = buf
	appFileSystem = memfs

	// Act
	rootCmd.SetArgs([]string{"lp", "-p", dir})
	rootCmd.Execute()

	// Assert
	actual := buf.String()
	ass.Equal(`
These projects are included into a solution but not found in the file system:

Solution: a\a.sln
 a\a\a.csproj
`, actual)
}
