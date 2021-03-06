package cmd

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"solt/internal/out"
	"solt/internal/sys"
	"testing"
)

func Test_FindLostFilesCmd_NoLostFilesFound(t *testing.T) {
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

	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lf", dir)

	// Assert
	actual := env.String()
	ass.Equal(`
 Totals:

  Projects                  1
                            
  Files                     Count
  -----                     -----
  Found                     2
  Included                  4
  Included but not exist    0
  Lost                      0
`, actual)
}

func Test_FindLostFilesCmdFilesInExcludedFolder_NoLostFilesFound(t *testing.T) {
	var tests = []struct {
		path string
	}{
		{"packages/Program.cs"},
		{"a/packages/Program.cs"},
		{"a/obj/Debug/Program.cs"},
	}
	for _, tst := range tests {
		// Arrange
		ass := assert.New(t)
		dir := "a/"
		memfs := afero.NewMemMapFs()
		_ = memfs.MkdirAll(dir+"a/Properties", 0755)
		_ = memfs.MkdirAll(dir+"packages", 0755)
		_ = memfs.MkdirAll(dir+"a/packages", 0755)
		_ = afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
		_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent), 0644)
		_ = afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
		_ = afero.WriteFile(memfs, dir+"a/packages.config", []byte(packagesConfingContent), 0644)
		_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
		_ = afero.WriteFile(memfs, dir+tst.path, []byte(codeFileContent), 0644)
		_ = afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

		env := out.NewMemoryEnvironment()

		// Act
		_ = Execute(memfs, env, "lf", dir)

		// Assert
		actual := env.String()
		ass.Equal(`
 Totals:

  Projects                  1
                            
  Files                     Count
  -----                     -----
  Found                     3
  Included                  4
  Included but not exist    0
  Lost                      0
`, actual)
	}
}

func Test_FindLostFilesCmd_LostFilesFound(t *testing.T) {
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
	_ = afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo1.cs", []byte(assemblyInfoContent), 0644)

	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lf", dir)

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`  a\a\Properties\AssemblyInfo1.cs

 Totals:

  Projects                  1
                            
  Files                     Count
  -----                     -----
  Found                     3
  Included                  4
  Included but not exist    0
  Lost                      1
`), actual)
}

func Test_FindLostFilesCmdNoProjects_AllFilesLost(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	_ = memfs.MkdirAll(dir+"a/Properties", 0755)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lf", dir)

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`  a\a\Program.cs
  a\a\Properties\AssemblyInfo.cs

 Totals:

  Projects                  0
                            
  Files                     Count
  -----                     -----
  Found                     2
  Included                  0
  Included but not exist    0
  Lost                      2
`), actual)
}

func Test_FindLostFilesCmdSeveralSolutions_LostFilesFound(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	root := "root/"
	memfs := afero.NewMemMapFs()

	dir1 := root + "a/"
	_ = memfs.MkdirAll(dir1+"a/Properties", 0755)
	_ = afero.WriteFile(memfs, dir1+"a.sln", []byte(testSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir1+"a/a.csproj", []byte(testProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir1+"a/App.config", []byte(appConfigContent), 0644)
	_ = afero.WriteFile(memfs, dir1+"a/packages.config", []byte(packagesConfingContent), 0644)
	_ = afero.WriteFile(memfs, dir1+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir1+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)
	_ = afero.WriteFile(memfs, dir1+"a/Properties/AssemblyInfo1.cs", []byte(assemblyInfoContent), 0644)

	dir2 := root + "b/"
	_ = memfs.MkdirAll(dir2+"a/Properties", 0755)
	_ = afero.WriteFile(memfs, dir2+"a.sln", []byte(testSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir2+"a/a.csproj", []byte(testProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir2+"a/App.config", []byte(appConfigContent), 0644)
	_ = afero.WriteFile(memfs, dir2+"a/packages.config", []byte(packagesConfingContent), 0644)
	_ = afero.WriteFile(memfs, dir2+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir2+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)
	_ = afero.WriteFile(memfs, dir2+"a/Properties/AssemblyInfo1.cs", []byte(assemblyInfoContent), 0644)

	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lf", root)

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`  root\a\a\Properties\AssemblyInfo1.cs
  root\b\a\Properties\AssemblyInfo1.cs

 Totals:

  Projects                  2
                            
  Files                     Count
  -----                     -----
  Found                     6
  Included                  8
  Included but not exist    0
  Lost                      2
`), actual)
}

func Test_FindLostFilesCmdSdkProjects_NoLostFilesFound(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte(coreSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(aSdkProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"b/b.csproj", []byte(bSdkProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"b/Class1.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"c/c.csproj", []byte(cSdkProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"c/Class1.cs", []byte(codeFileContent), 0644)

	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lf", dir)

	// Assert
	actual := env.String()
	ass.Equal(`
 Totals:

  Projects                  3
                            
  Files                     Count
  -----                     -----
  Found                     3
  Included                  0
  Included but not exist    0
  Lost                      0
`, actual)
}

func Test_FindLostFilesCmdExplicitFilterSet_LostFilesFound(t *testing.T) {
	var tests = []struct {
		filter string
	}{
		{".cs"},
		{".CS"},
	}

	for _, tst := range tests {
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
		_ = afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo1.cs", []byte(assemblyInfoContent), 0644)

		env := out.NewMemoryEnvironment()

		// Act
		_ = Execute(memfs, env, "lf", dir, "-f", tst.filter)

		// Assert
		actual := env.String()
		ass.Equal(sys.ToValidPath(`  a\a\Properties\AssemblyInfo1.cs

 Totals:

  Projects                  1
                            
  Files                     Count
  -----                     -----
  Found                     3
  Included                  4
  Included but not exist    0
  Lost                      1
`), actual)
	}
}

func Test_FindLostFilesCmdRemove_LostFilesRemoved(t *testing.T) {
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
	_ = afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo1.cs", []byte(assemblyInfoContent), 0644)

	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lf", dir, "-r")

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`  a\a\Properties\AssemblyInfo1.cs
File: a\a\Properties\AssemblyInfo1.cs removed successfully.

 Totals:

  Projects                  1
                            
  Files                     Count
  -----                     -----
  Found                     3
  Included                  4
  Included but not exist    0
  Lost                      1
`), actual)
	_, err := memfs.Stat(dir + "a/Properties/AssemblyInfo1.cs")
	ass.Error(err)
}

func Test_FindLostFilesCmdRemoveReadOnly_LostFilesNotRemoved(t *testing.T) {
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
	_ = afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo1.cs", []byte(assemblyInfoContent), 0644)

	fs := afero.NewReadOnlyFs(memfs)
	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(fs, env, "lf", dir, "-r")

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`  a\a\Properties\AssemblyInfo1.cs

 Totals:

  Projects                  1
                            
  Files                     Count
  -----                     -----
  Found                     3
  Included                  4
  Included but not exist    0
  Lost                      1
`), actual)
	_, err := fs.Stat(dir + "a/Properties/AssemblyInfo1.cs")
	ass.NoError(err)
}

func Test_FindLostFilesCmdUnexistOptionEnabled_UnesistFilesFound(t *testing.T) {
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

	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lf", "-a", dir)

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`
These files included into projects but not exist in the file system.

 Project: a\a\a.csproj
  a\a\Properties\AssemblyInfo.cs

 Totals:

  Projects                  1
                            
  Files                     Count
  -----                     -----
  Found                     1
  Included                  4
  Included but not exist    1
  Lost                      0
`), actual)
}

func Test_FindLostFilesCmdUnexistOptionNotSet_UnesistFilesNotShown(t *testing.T) {
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

	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lf", dir)

	// Assert
	actual := env.String()
	ass.Equal(`
 Totals:

  Projects                  1
                            
  Files                     Count
  -----                     -----
  Found                     1
  Included                  4
  Included but not exist    0
  Lost                      0
`, actual)
}

func Test_FindLostFilesNoPath_OutputHelp(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lf")

	// Assert
	actual := env.String()
	ass.Contains(actual, "Find lost files in the folder specified")
}

func Test_FindLostFilesEmptyPath_NoOutput(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lf", "/")

	// Assert
	actual := env.String()
	ass.Equal(`
 Totals:

  Projects                  0
                            
  Files                     Count
  -----                     -----
  Found                     0
  Included                  0
  Included but not exist    0
  Lost                      0
`, actual)
}
