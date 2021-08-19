package cmd

import (
	"bytes"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"io"
	"solt/internal/out"
	"solt/internal/sys"
	"strings"
	"testing"
)

func Test_ValidateSdkSolutionCmd_RedundantReferencesFound(t *testing.T) {
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
	_ = Execute(memfs, env, "va", dir)

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`
 Solution: a\a.sln
   project a\a\a.csproj has redundant references:
     a\b\b.csproj


 Totals:

  Parameter               Count    %     
  ---------               -----    ------
  Solutions               1        
  Problem solutions       1        100.00%
  SDK Projects            3        
  Problem projects        1        33.33%
  Redundant references    1
`), actual)
}

func Test_ValidateSdkSolutionCmdSeveralSolutionsWithSameProjects_RedundantReferencesFound(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	top := "t/"
	dir := top + "a/"
	memfs := afero.NewMemMapFs()
	_ = afero.WriteFile(memfs, top+"t.sln", []byte(coreTopSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte(coreSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(aSdkProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"b/b.csproj", []byte(bSdkProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"b/Class1.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"c/c.csproj", []byte(cSdkProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"c/Class1.cs", []byte(codeFileContent), 0644)

	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "va", top)

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(`
 Solution: t\a\a.sln
   project t\a\a\a.csproj has redundant references:
     t\a\b\b.csproj


 Solution: t\t.sln
   project t\a\a\a.csproj has redundant references:
     t\a\b\b.csproj


 Totals:

  Parameter               Count    %     
  ---------               -----    ------
  Solutions               2        
  Problem solutions       2        100.00%
  SDK Projects            3        
  Problem projects        1        33.33%
  Redundant references    1
`), actual)
}

func Test_FixSdkSolutionCmd_RedundantReferencesRemoved(t *testing.T) {
	var tests = []struct {
		name      string
		redundant string
		expect    string
	}{
		{"unix", aSdkProjectContent, aSdkProjectContentWithoutRedundantRefs},
		{"unix full tags", aSdkProjectContentFullTags, aSdkProjectContentFullTagsNoRedundant},
		{"windows", unix2Win(aSdkProjectContent), unix2Win(aSdkProjectContentWithoutRedundantRefs)},
		{"windows full tags", unix2Win(aSdkProjectContentFullTags), unix2Win(aSdkProjectContentFullTagsNoRedundant)},
	}
	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)
			dir := "a/"
			memfs := afero.NewMemMapFs()

			_ = afero.WriteFile(memfs, dir+"a.sln", []byte(coreSolutionContent), 0644)
			_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(tst.redundant), 0644)
			_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
			_ = afero.WriteFile(memfs, dir+"b/b.csproj", []byte(bSdkProjectContent), 0644)
			_ = afero.WriteFile(memfs, dir+"b/Class1.cs", []byte(codeFileContent), 0644)
			_ = afero.WriteFile(memfs, dir+"c/c.csproj", []byte(cSdkProjectContent), 0644)
			_ = afero.WriteFile(memfs, dir+"c/Class1.cs", []byte(codeFileContent), 0644)

			env := out.NewMemoryEnvironment()

			// Act
			_ = Execute(memfs, env, "va", "fix", dir)

			// Assert
			actual := env.String()
			ass.Equal(sys.ToValidPath("Fixed 1 redundant project references in 1 projects within solution a\\a.sln\n"), actual)
			fa, _ := memfs.Open(dir + "a/a.csproj")
			buf := bytes.NewBuffer(nil)
			_, _ = io.Copy(buf, fa)
			ass.Equal(tst.expect, string(buf.Bytes()))
		})
	}
}

func Test_FixSdkSolutionCmdSeveralSolutionsCase_RedundantReferencesRemoved(t *testing.T) {
	var tests = []struct {
		name      string
		redundant string
		expect    string
	}{
		{"unix", aSdkProjectContent, aSdkProjectContentWithoutRedundantRefs},
		{"unix full tags", aSdkProjectContentFullTags, aSdkProjectContentFullTagsNoRedundant},
		{"windows", unix2Win(aSdkProjectContent), unix2Win(aSdkProjectContentWithoutRedundantRefs)},
		{"windows full tags", unix2Win(aSdkProjectContentFullTags), unix2Win(aSdkProjectContentFullTagsNoRedundant)},
	}
	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)
			top := "t/"
			dir := top + "a/"
			memfs := afero.NewMemMapFs()
			_ = afero.WriteFile(memfs, top+"t.sln", []byte(coreTopSolutionContent), 0644)
			_ = afero.WriteFile(memfs, dir+"a.sln", []byte(coreSolutionContent), 0644)
			_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(tst.redundant), 0644)
			_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
			_ = afero.WriteFile(memfs, dir+"b/b.csproj", []byte(bSdkProjectContent), 0644)
			_ = afero.WriteFile(memfs, dir+"b/Class1.cs", []byte(codeFileContent), 0644)
			_ = afero.WriteFile(memfs, dir+"c/c.csproj", []byte(cSdkProjectContent), 0644)
			_ = afero.WriteFile(memfs, dir+"c/Class1.cs", []byte(codeFileContent), 0644)

			env := out.NewMemoryEnvironment()

			// Act
			_ = Execute(memfs, env, "va", "fix", top)

			// Assert
			actual := env.String()
			ass.Equal(sys.ToValidPath("Fixed 1 redundant project references in 1 projects within solution t\\a\\a.sln\nFixed 1 redundant project references in 1 projects within solution t\\t.sln\n"), actual)
			fa, _ := memfs.Open(dir + "a/a.csproj")
			buf := bytes.NewBuffer(nil)
			_, _ = io.Copy(buf, fa)
			ass.Equal(tst.expect, string(buf.Bytes()))
		})
	}
}

func unix2Win(s string) string {
	return strings.ReplaceAll(s, "\n", "\r\n")
}

func Test_ValidateOldSolutionCmd_RedundantReferencesNotFound(t *testing.T) {
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
	_ = Execute(memfs, env, "va", dir)

	// Assert
	actual := env.String()
	ass.Equal(`
 Totals:

  Parameter               Count    %     
  ---------               -----    ------
  Solutions               1        
  Problem solutions       0        0.00%
  SDK Projects            0        
  Problem projects        0        0.00%
  Redundant references    0
`, actual)
}

func Test_FixSdkSolutionCmd_RedundantReferencesNotFound(t *testing.T) {
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
	_ = Execute(memfs, env, "fr", dir)

	// Assert
	actual := env.String()
	ass.Equal("", actual)
}

func Test_ValidateSdkNoPath_OutputHelp(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "va")

	// Assert
	actual := env.String()
	ass.Contains(actual, "Validates SDK projects within solution(s)")
}

func Test_FixSdkNoPath_OutputHelp(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "va", "fix")

	// Assert
	actual := env.String()
	ass.Contains(actual, "Fixes redundant SDK projects references")
}
