package cmd

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ValidateSdkSolutionCmd_RedundantReferencesFound(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	afero.WriteFile(memfs, dir+"a.sln", []byte(coreSolutionContent), 0644)
	afero.WriteFile(memfs, dir+"a/a.csproj", []byte(aSdkProjectContent), 0644)
	afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	afero.WriteFile(memfs, dir+"b/b.csproj", []byte(bSdkProjectContent), 0644)
	afero.WriteFile(memfs, dir+"b/Class1.cs", []byte(codeFileContent), 0644)
	afero.WriteFile(memfs, dir+"c/c.csproj", []byte(cSdkProjectContent), 0644)
	afero.WriteFile(memfs, dir+"c/Class1.cs", []byte(codeFileContent), 0644)

	appPrinter = newMockPrn()
	appFileSystem = memfs

	// Act
	Execute("va", "-p", dir)

	// Assert
	actual := appPrinter.(*mockprn).String()
	ass.Equal(` Solution: <green>a\a.sln</>
   project: <bold>a\a\a.csproj</> has redundant references
    <gray>a\b\b.csproj</>
`, actual)
}

func Test_ValidateOldSolutionCmd_RedundantReferencesNotFound(t *testing.T) {
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
	Execute("va", "-p", dir)

	// Assert
	actual := appPrinter.(*mockprn).String()
	ass.Equal("", actual)
}
