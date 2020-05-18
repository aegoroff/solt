package cmd

import (
	"bytes"
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_FindLostFilesCmd_NoLostFilesFound(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	memfs.MkdirAll(dir + "a/Properties", 0755)
	afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent), 0644)
	afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)
	appFileSystem = memfs

	buf := bytes.NewBufferString("")
	appWriter = buf

	// Act
	rootCmd.SetArgs([]string{"lf", "-p", dir})
	rootCmd.Execute()

	// Assert
	actual := buf.String();
	ass.Equal(``, actual)
}

func Test_FindLostFilesCmd_LostFilesFound(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	memfs.MkdirAll(dir + "a/Properties", 0755)
	afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent), 0644)
	afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)
	afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo1.cs", []byte(assemblyInfoContent), 0644)
	appFileSystem = memfs

	buf := bytes.NewBufferString("")
	appWriter = buf

	// Act
	rootCmd.SetArgs([]string{"lf", "-p", dir})
	rootCmd.Execute()

	// Assert
	actual := buf.String();
	ass.Equal(" a\\a\\Properties\\AssemblyInfo1.cs\n", actual)
}

func Test_FindLostFiles(t *testing.T) {
	// Arrange
	f1 := `c:\prj\f1\p1.csproj`
	f2 := `c:\prj\f2\p2.csproj`
	f3 := `c:\prj\f3\p3.csproj`

	p1 := Project{
		OutputPaths: []string{`bin\Debug`, `bin\Release`},
		Compiles:    []Include{{Path: "code1.cs"}, {Path: "code2.cs"}},
	}

	p2 := Project{
		OutputPaths: []string{`bin\Debug`, `bin\Release`},
		Compiles:    []Include{{Path: "code1.cs"}, {Path: "code2.cs"}},
	}

	fi1 := folderInfo{
		projectPath: &f1,
		project:     &p1,
	}
	fi2 := folderInfo{
		projectPath: &f2,
		project:     &p2,
	}
	fi3 := folderInfo{
		projectPath: &f3,
	}

	infos := []*folderInfo{&fi1, &fi2, &fi3}

	tree := rbtree.NewRbTree()

	for i, fi := range infos {
		path := fmt.Sprintf(`c:\prj\f%d\`, i+1)
		key := newProjectTreeNode(path, fi)
		n := rbtree.NewNode(key)
		tree.Insert(n)
	}

	ass := assert.New(t)
	var tests = []struct {
		foundfiles []string
		result     []string
	}{
		{[]string{`c:\prj\f1\code1.cs`, `c:\prj\f1\code2.cs`, `c:\\prj\f1\code3.cs`, `c:\prj\f2\code1.cs`, `c:\prj\f2\code2.cs`, `c:\prj\f2\code3.cs`}, []string{`c:\\prj\f1\code3.cs`, `c:\prj\f2\code3.cs`}},
		{[]string{`c:\prj\f1\code1.cs`, `c:\prj\f1\code2.cs`, `c:\prj\f2\code1.cs`, `c:\prj\f2\code2.cs`}, []string(nil)},
		{[]string{`c:\prj\f1\cOde1.cs`, `c:\prj\f1\Code2.cs`, `c:\prj\f2\coDe1.cs`, `c:\prj\f2\Code2.cs`}, []string(nil)},
		{[]string{`c:\prj\f1\code1.cs`, `c:\prj\f1\code2.cs`, `c:\prj\f1\bin\Debug\code3.cs`, `c:\prj\f2\code1.cs`, `c:\prj\f2\code2.cs`, `c:\prj\f2\code3.cs`}, []string{`c:\prj\f2\code3.cs`}},
		{[]string{`c:\prj\f1\code1.cs`, `c:\prj\f1\code2.cs`, `c:\prj\f1\bin\Release\code3.cs`, `c:\prj\f2\code1.cs`, `c:\prj\f2\code2.cs`, `c:\prj\f2\code3.cs`}, []string{`c:\prj\f2\code3.cs`}},
		{[]string{`c:\prj\f1\cOde1.cs`, `c:\prj\f1\Code2.cs`, `c:\prj\f2\coDe1.cs`, `c:\prj\f2\Code2.cs`, `c:\prj\f3\Code1.cs`}, []string{`c:\prj\f3\Code1.cs`}},
	}

	for _, test := range tests {
		// Act
		result, unexists := findLostFiles(tree, afero.NewMemMapFs(), map[string]interface{}{`c:\prj\packages`: nil}, test.foundfiles)

		// Assert
		ass.Equal(test.result, result)
		ass.Equal(2, len(unexists))
	}
}
