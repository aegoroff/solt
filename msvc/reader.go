package msvc

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/afero"
	"solt/internal/sys"
	"sync"
)

// ReadSolutionDir reads filesystem directory and all its childs to get information
// about all solutions and projects in this tree.
// It returns tree
func ReadSolutionDir(path string, fs afero.Fs, fileHandlers ...ReaderHandler) rbtree.RbTree {
	result := rbtree.NewRbTree()

	aggregateChannel := make(chan *Folder, 4)
	fileChannel := make(chan string, 16)

	var wg sync.WaitGroup

	// Aggregating goroutine
	go func() {
		defer wg.Done()
		for folder := range aggregateChannel {
			current, ok := result.Search(folder)
			if !ok {
				// Create new node
				result.Insert(folder)
			} else {
				// Update current, that has already been created before
				folder.copyContent(current.Key().(*Folder))
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
