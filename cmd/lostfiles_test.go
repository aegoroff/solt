package cmd

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"solt/cmd/api"
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

	env := api.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lf", "-p", dir)

	// Assert
	actual := env.String()
	ass.Equal(``, actual)
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

		env := api.NewMemoryEnvironment()

		// Act
		_ = Execute(memfs, env, "lf", "-p", dir)

		// Assert
		actual := env.String()
		ass.Equal(``, actual)
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

	env := api.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lf", "-p", dir)

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath("  a\\a\\Properties\\AssemblyInfo1.cs\n"), actual)
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

	env := api.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lf", "-p", root)

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath("  root\\a\\a\\Properties\\AssemblyInfo1.cs\n  root\\b\\a\\Properties\\AssemblyInfo1.cs\n"), actual)
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

	env := api.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lf", "-p", dir)

	// Assert
	actual := env.String()
	ass.Equal("", actual)
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

		env := api.NewMemoryEnvironment()

		// Act
		_ = Execute(memfs, env, "lf", "-p", dir, "-f", tst.filter)

		// Assert
		actual := env.String()
		ass.Equal(sys.ToValidPath("  a\\a\\Properties\\AssemblyInfo1.cs\n"), actual)
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

	env := api.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lf", "-p", dir, "-r")

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath("  a\\a\\Properties\\AssemblyInfo1.cs\nFile: a\\a\\Properties\\AssemblyInfo1.cs removed successfully.\n"), actual)
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
	env := api.NewMemoryEnvironment()

	// Act
	_ = Execute(fs, env, "lf", "-p", dir, "-r")

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath("  a\\a\\Properties\\AssemblyInfo1.cs\n"), actual)
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

	env := api.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lf", "-p", dir, "-a")

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath("\nThese files included into projects but not exist in the file system.\n\n Project: a\\a\\a.csproj\n  a\\a\\Properties\\AssemblyInfo.cs\n"), actual)
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

	env := api.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lf", "-p", dir)

	// Assert
	actual := env.String()
	ass.Equal("", actual)
}

func Test_FindLostFilesNoPath_OutputHelp(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	env := api.NewMemoryEnvironment()

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
	env := api.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "lf", "-p", "/")

	// Assert
	actual := env.String()
	ass.Equal("", actual)
}
