package msvc

import (
	"github.com/google/btree"
	"github.com/spf13/afero"
	"solt/internal/sys"
	"sync"
)

// ReadSolutionDir reads filesystem directory and all its childs to get information
// about all solutions and projects in this tree.
// It returns tree
func ReadSolutionDir(path string, fs afero.Fs, fileHandlers ...ReaderHandler) *btree.BTree {
	result := btree.New(16)

	aggregateChannel := make(chan *Folder, 4)
	fileChannel := make(chan string, 16)

	var wg sync.WaitGroup

	// Aggregating goroutine
	go func() {
		defer wg.Done()
		for f := range aggregateChannel {
			current := result.Get(f)
			if current == nil {
				// Create new node
				result.ReplaceOrInsert(f)
			} else {
				// Update folder node that has already been created before
				merge(current.(*Folder), f)
			}
		}
	}()

	modules := newReaderModules(fs)

	rdr := reader{aggregator: aggregateChannel, modules: modules}

	fhandlers := []ReaderHandler{&rdr}
	fhandlers = append(fhandlers, fileHandlers...)

	// Reading files goroutine
	go func(handlers []ReaderHandler) {
		defer close(aggregateChannel)

		for path := range fileChannel {
			for _, h := range handlers {
				h.Handler(path)
			}
		}
	}(fhandlers)

	handlers := []sys.ScanHandler{func(evt *sys.ScanEvent) {
		if evt.File == nil {
			return
		}
		f := evt.File
		fileChannel <- f.Path
	}}

	// Start reading path
	wg.Add(1)

	sys.Scan(path, fs, handlers)

	close(fileChannel)

	wg.Wait()

	return result
}

func merge(to, from *Folder) {
	toC := to.Content
	fromC := from.Content
	if fromC.Packages != nil {
		toC.Packages = fromC.Packages
	} else {
		toC.Projects = append(toC.Projects, fromC.Projects...)
		toC.Solutions = append(toC.Solutions, fromC.Solutions...)
	}
}
