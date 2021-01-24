package msvc

import "solt/solution"

// ProjectIterator provides iteration over all solution projects
// eacj of them is *MsbuildProject
type ProjectIterator struct {
	sln    *VisualStudioSolution
	search ProjectSearcher
}

// NewProjectIterator creates new ProjectIterator instance
func NewProjectIterator(sln *VisualStudioSolution, search ProjectSearcher) *ProjectIterator {
	return &ProjectIterator{sln: sln, search: search}
}

// Foreach goes through all solution's projects converts them into *MsbuildProject
// and calls function specified
func (pi *ProjectIterator) Foreach(callFn func(*MsbuildProject)) {
	pi.sln.Projects(func(project *solution.Project) {
		found, ok := pi.search.Search(project)
		if ok {
			callFn(found)
		}
	})
}
