package cmd

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"solt/cmd/out"
	"solt/internal/sys"
	"testing"
)

func Test_FindLostProjectsCmd_NoLostProjectsFound(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()

	_ = memfs.MkdirAll(dir+"a/Properties", 0755)
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lp", "-p", dir)

	// Assert
	actual := env.String()
	ass.Equal(``, actual)
}

func Test_FindLostProjectsCmdLostProjectsInTheSameDir_LostProjectsFound(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	_ = memfs.MkdirAll(dir+"a/Properties", 0755)
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/a1.csproj", []byte(testProjectContent2), 0644)
	_ = afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/packages.config", []byte(packagesConfingContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lp", "-p", dir)

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`
 These projects are not included into any solution
 but files from the projects' folders are used in another projects within a solution:

  a\a\a1.csproj
`), actual)
}

func Test_FindLostProjectsCmdLostProjectsRemove_LostProjectsRemoved(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	_ = memfs.MkdirAll(dir+"a/Properties", 0755)
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a1/a1.csproj", []byte(testProjectContent3), 0644)
	_ = afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a1/App.config", []byte(appConfigContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/packages.config", []byte(packagesConfingContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lp", "-p", dir, "-r")

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`  a\a1\a1.csproj

 Folder 'a\a1\' removed
`), actual)
}

func Test_FindLostProjectsCmdLostProjectsRemoveReadOnlyFs_LostProjectsNotRemoved(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	_ = memfs.MkdirAll(dir+"a/Properties", 0755)
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a1/a1.csproj", []byte(testProjectContent3), 0644)
	_ = afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a1/App.config", []byte(appConfigContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/packages.config", []byte(packagesConfingContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(afero.NewReadOnlyFs(memfs), env, "lp", "-p", dir, "-r")

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`  a\a1\a1.csproj
`), actual)
}

func Test_FindLostProjectsCmdOtherDirWithFilesIncludedToLinked_LostProjectsFound(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	_ = memfs.MkdirAll(dir+"a/Properties", 0755)
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent3a1), 0644)
	_ = afero.WriteFile(memfs, dir+"a1/a1.csproj", []byte(testProjectContent3), 0644)
	_ = afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a1/App.config", []byte(appConfigContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/packages.config", []byte(packagesConfingContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a1/packages.config", []byte(packagesConfingContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a1/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a1/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lp", "-p", dir)

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`
 These projects are not included into any solution
 but files from the projects' folders are used in another projects within a solution:

  a\a1\a1.csproj
`), actual)
}

func Test_FindLostProjectsCmdOtherDirWithFilesDeepIncludedToLinked_LostProjectsFound(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	_ = memfs.MkdirAll(dir+"a/Properties", 0755)
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent3a1b), 0644)
	_ = afero.WriteFile(memfs, dir+"a1/a1.csproj", []byte(testProjectContent3), 0644)
	_ = afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a1/App.config", []byte(appConfigContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/packages.config", []byte(packagesConfingContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a1/packages.config", []byte(packagesConfingContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a1/b/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a1/b/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lp", "-p", dir)

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`
 These projects are not included into any solution
 but files from the projects' folders are used in another projects within a solution:

  a\a1\a1.csproj
`), actual)
}

func Test_FindLostProjectsCmdLostProjectsInOtherDir_LostProjectsFound(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	_ = memfs.MkdirAll(dir+"a/Properties", 0755)
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a1/a1.csproj", []byte(testProjectContent2), 0644)
	_ = afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/packages.config", []byte(packagesConfingContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lp", "-p", dir)

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`  a\a1\a1.csproj
`), actual)
}

func Test_FindLostProjectsCmdUnexistProjects_LostProjectsFound(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	_ = memfs.MkdirAll(dir+"a/Properties", 0755)
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/packages.config", []byte(packagesConfingContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lp", "-p", dir)

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`
These projects are included into a solution but not found in the file system:

 Solution: a\a.sln
  a\a\a.csproj
`), actual)
}

func Test_FindLostProjectsNoPath_OutputHelp(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lp")

	// Assert
	actual := env.String()
	ass.Contains(actual, "Find projects that not included into any solution")
}
