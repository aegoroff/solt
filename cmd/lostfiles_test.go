package cmd

import (
	"bytes"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_FindLostFilesCmd_NoLostFilesFound(t *testing.T) {
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
	rootCmd.SetArgs([]string{"lf", "-p", dir})
	rootCmd.Execute()

	// Assert
	actual := buf.String()
	ass.Equal(``, actual)
}

func Test_FindLostFilesCmd_LostFilesFound(t *testing.T) {
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
	afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo1.cs", []byte(assemblyInfoContent), 0644)

	buf := bytes.NewBufferString("")

	appWriter = buf
	appFileSystem = memfs

	// Act
	rootCmd.SetArgs([]string{"lf", "-p", dir})
	rootCmd.Execute()

	// Assert
	actual := buf.String()
	ass.Equal(" a\\a\\Properties\\AssemblyInfo1.cs\n", actual)
}

func Test_FindLostFilesCmd_UnesistFilesFound(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	memfs.MkdirAll(dir+"a/Properties", 0755)
	afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent), 0644)
	afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)

	buf := bytes.NewBufferString("")

	appWriter = buf
	appFileSystem = memfs

	// Act
	rootCmd.SetArgs([]string{"lf", "-p", dir})
	rootCmd.Execute()

	// Assert
	actual := buf.String()
	ass.Equal("\nThese files included into projects but not exist in the file system.\n\nProject: a\\a\\a.csproj\n a\\a\\Properties\\AssemblyInfo.cs\n", actual)
}
