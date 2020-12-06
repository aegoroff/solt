package msvc

import (
	"github.com/aegoroff/godatastruct/collections"
	"github.com/aegoroff/godatastruct/rbtree"
	"path/filepath"
	"solt/solution"
)

// SelectAllSolutionProjectPaths gets all possible projects' paths defined in solution
func SelectAllSolutionProjectPaths(sln *VisualStudioSolution, pathDecorator StringDecorator) collections.StringHashSet {
	solutionPath := filepath.Dir(sln.Path)
	var paths = make(collections.StringHashSet, len(sln.Solution.Projects))
	for _, sp := range sln.Solution.Projects {
		if sp.TypeID == solution.IDSolutionFolder {
			continue
		}
		fullProjectPath := filepath.Join(solutionPath, sp.Path)
		paths.Add(pathDecorator(fullProjectPath))
	}
	return paths
}

// GetFilesIncludedIntoProject gets all files included into MSBuild project
func GetFilesIncludedIntoProject(prj *MsbuildProject) []string {
	var result []string
	folderPath := filepath.Dir(prj.Path)

	msp := prj.Project

	l := len(msp.Contents) + len(msp.Nones) + len(msp.CLCompiles) + len(msp.CLInclude) + len(msp.Compiles)
	includes := make([]include, 0, l)
	includes = append(includes, msp.Contents...)
	includes = append(includes, msp.Nones...)
	includes = append(includes, msp.CLCompiles...)
	includes = append(includes, msp.CLInclude...)
	includes = append(includes, msp.Compiles...)

	result = append(result, createPathsFromIncludes(includes, folderPath)...)

	return result
}

// SelectProjects gets all Visual Studion solutions found in a directory
func SelectProjects(foldersTree rbtree.RbTree) []*MsbuildProject {
	var projects []*MsbuildProject

	WalkProjectFolders(foldersTree, func(project *MsbuildProject, folder *Folder) {
		projects = append(projects, project)
	})

	return projects
}

// SelectSolutions gets all Visual Studion solutions found in a directory
func SelectSolutions(foldersTree rbtree.RbTree) []*VisualStudioSolution {
	w := walkSol{solutions: make([]*VisualStudioSolution, 0)}
	walk(foldersTree, &w)
	return w.solutions
}

// SelectSolutionsAndProjects gets all Visual Studion solutions and projects found in a directory
func SelectSolutionsAndProjects(foldersTree rbtree.RbTree) ([]*VisualStudioSolution, []*MsbuildProject) {
	ws := walkSol{solutions: make([]*VisualStudioSolution, 0)}
	var projects []*MsbuildProject

	wp := walkPrj{handler: func(project *MsbuildProject, folder *Folder) {
		projects = append(projects, project)
	}}

	walk(foldersTree, &ws, &wp)
	return ws.solutions, projects
}

// WalkProjectFolders traverse all projects found in solution(s) folder
func WalkProjectFolders(foldersTree rbtree.RbTree, action ProjectHandler) {
	w := &walkPrj{handler: action}
	walk(foldersTree, w)
}

// IsSdkProject gets whether a project is a the new VS 2017 or later project
func (p *msbuildProject) IsSdkProject() bool {
	if len(p.Sdk) > 0 {
		return true
	}
	if len(p.Imports) == 0 {
		return false
	}
	for _, imp := range p.Imports {
		if len(imp.Sdk) > 0 {
			return true
		}
	}
	return false
}

func createPathsFromIncludes(paths []include, basePath string) []string {
	if paths == nil {
		return []string{}
	}

	result := make([]string, 0, len(paths))

	for _, c := range paths {
		fp := filepath.Join(basePath, c.Path)
		result = append(result, fp)
	}

	return result
}
