package msvc

import "solt/solution"

// ReaderHandler defines file system scanning handler
type ReaderHandler interface {
	// Handler method called on each file and folder scanned
	Handler(path string)
}

type readerModule interface {
	allow(path string) bool
	read(path string, ch chan<- *Folder)
}

type walker interface {
	walk(f *Folder)
}

// ProjectSearcher provides *MsbuildProject searching using
// *solution.Project
type ProjectSearcher interface {
	// Search search *MsbuildProject
	Search(*solution.Project) (*MsbuildProject, bool)
}
