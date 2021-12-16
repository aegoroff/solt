package cmd

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"solt/internal/out"
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
	_ = Execute(memfs, env, "lp", dir)

	// Assert
	actual := env.String()
	ass.Equal(`
 Totals:

  Solutions                      1
  Projects                       1
                                 
                                 Count    %     
                                 -----    ------
  Within solutions               1        100.00%
  Lost projects                  0        0.00%
  Lost projects with includes    0        0.00%
  Included but not exist         0        0.00%
  Removed (if specified)         0        0.00%
`, actual)
}

func Test_FindLostProjectsCmdPureSdkProjects_LostProjectsFound(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()

	_ = memfs.MkdirAll(dir+"a/Properties", 0755)
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte(coreSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(aSdkProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"b/b.csproj", []byte(bSdkProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"b1/b.csproj", []byte(bSdkProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"b/Class1.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"c/c.csproj", []byte(cSdkProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"c/Class1.cs", []byte(codeFileContent), 0644)

	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lp", dir)

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`  a\b1\b.csproj

 Totals:

  Solutions                      1
  Projects                       4
                                 
                                 Count    %     
                                 -----    ------
  Within solutions               3        75.00%
  Lost projects                  1        25.00%
  Lost projects with includes    0        0.00%
  Included but not exist         0        0.00%
  Removed (if specified)         0        0.00%
`), actual)
}

func Test_FindLostProjectsCmdMissingProjectInSeveralSolutions_LostProjectsFound(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()

	_ = memfs.MkdirAll(dir+"a/Properties", 0755)
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte(coreSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a1.sln", []byte(coreSolutionContent2), 0644)
	_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(aSdkProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a1/a.csproj", []byte(aSdkProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a1/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"b/b.csproj", []byte(bSdkProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"b/Class1.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"b1/b.csproj", []byte(bSdkProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"b1/Class1.cs", []byte(codeFileContent), 0644)

	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lp", dir)

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`
These projects are included into a solution but not found in the file system:

 Solution: a\a.sln
  a\c\c.csproj

 Solution: a/a1.sln
  a\c\c.csproj

 Totals:

  Solutions                      2
  Projects                       4
                                 
                                 Count    %     
                                 -----    ------
  Within solutions               4        100.00%
  Lost projects                  0        0.00%
  Lost projects with includes    0        0.00%
  Included but not exist         1        25.00%
  Removed (if specified)         0        0.00%
`), actual)
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
	_ = Execute(memfs, env, "lp", dir)

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`
 These projects are not included into any solution
 but files from the projects' folders are used in another projects within a solution:

  a\a\a1.csproj

 Totals:

  Solutions                      1
  Projects                       2
                                 
                                 Count    %     
                                 -----    ------
  Within solutions               1        50.00%
  Lost projects                  0        0.00%
  Lost projects with includes    1        50.00%
  Included but not exist         0        0.00%
  Removed (if specified)         0        0.00%
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
	_ = Execute(memfs, env, "lp", dir, "-r")

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`  a\a1\a1.csproj

 Folder 'a\a1\' removed

 Totals:

  Solutions                      1
  Projects                       2
                                 
                                 Count    %     
                                 -----    ------
  Within solutions               1        50.00%
  Lost projects                  1        50.00%
  Lost projects with includes    0        0.00%
  Included but not exist         0        0.00%
  Removed (if specified)         1        50.00%
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
	_ = Execute(afero.NewReadOnlyFs(memfs), env, "lp", dir, "-r")

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`  a\a1\a1.csproj

 Totals:

  Solutions                      1
  Projects                       2
                                 
                                 Count    %     
                                 -----    ------
  Within solutions               1        50.00%
  Lost projects                  1        50.00%
  Lost projects with includes    0        0.00%
  Included but not exist         0        0.00%
  Removed (if specified)         0        0.00%
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
	_ = Execute(memfs, env, "lp", dir)

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`
 These projects are not included into any solution
 but files from the projects' folders are used in another projects within a solution:

  a\a1\a1.csproj

 Totals:

  Solutions                      1
  Projects                       2
                                 
                                 Count    %     
                                 -----    ------
  Within solutions               1        50.00%
  Lost projects                  0        0.00%
  Lost projects with includes    1        50.00%
  Included but not exist         0        0.00%
  Removed (if specified)         0        0.00%
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
	_ = Execute(memfs, env, "lp", dir)

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`
 These projects are not included into any solution
 but files from the projects' folders are used in another projects within a solution:

  a\a1\a1.csproj

 Totals:

  Solutions                      1
  Projects                       2
                                 
                                 Count    %     
                                 -----    ------
  Within solutions               1        50.00%
  Lost projects                  0        0.00%
  Lost projects with includes    1        50.00%
  Included but not exist         0        0.00%
  Removed (if specified)         0        0.00%
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
	_ = afero.WriteFile(memfs, dir+"a1/App.config", []byte(appConfigContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/packages.config", []byte(packagesConfingContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a1/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a1/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lp", dir)

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`  a\a1\a1.csproj

 Totals:

  Solutions                      1
  Projects                       2
                                 
                                 Count    %     
                                 -----    ------
  Within solutions               1        50.00%
  Lost projects                  1        50.00%
  Lost projects with includes    0        0.00%
  Included but not exist         0        0.00%
  Removed (if specified)         0        0.00%
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
	_ = Execute(memfs, env, "lp", dir)

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`
These projects are included into a solution but not found in the file system:

 Solution: a\a.sln
  a\a\a.csproj

 Totals:

  Solutions                      1
  Projects                       0
                                 
                                 Count    %     
                                 -----    ------
  Within solutions               0        0.00%
  Lost projects                  0        0.00%
  Lost projects with includes    0        0.00%
  Included but not exist         1        0.00%
  Removed (if specified)         0        0.00%
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
