package msvc

import (
	"github.com/aegoroff/dirstat/scan"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/afero"
	"sync"
)

// ReadSolutionDir reads filesystem directory and all its childs to get information
// about all solutions and projects in this tree.
// It returns tree
func ReadSolutionDir(path string, fs afero.Fs, fileHandlers ...ReaderHandler) rbtree.RbTree {
	modules := newReaderModules(fs)

	rdr := newReader(modules)

	fh := newFileEventHanlder(rdr)
	fh.addHandlers(fileHandlers...)

	go rdr.aggregate()

	scan.Scan(path, newFs(fs), fh)

	return rdr.getResult()
}

type fileEventHanlder struct {
	handlers []ReaderHandler
}

func newFileEventHanlder(rdr ReaderHandler) *fileEventHanlder {
	return &fileEventHanlder{
		handlers: []ReaderHandler{rdr},
	}
}

func (f *fileEventHanlder) addHandlers(handlers ...ReaderHandler) {
	f.handlers = append(f.handlers, handlers...)
}

func (f *fileEventHanlder) Handle(evt *scan.Event) {
	if evt.File != nil {
		for _, h := range f.handlers {
			h.Handler(evt.File.Path)
		}
	}
}

type reader struct {
	modules    []readerModule
	aggregator chan *Folder
	result     rbtree.RbTree
	wg         sync.WaitGroup
}

func (r *reader) getResult() rbtree.RbTree {
	close(r.aggregator)
	r.wg.Wait()
	return r.result
}

func newReader(modules []readerModule) *reader {
	return &reader{
		modules:    modules,
		result:     rbtree.NewRbTree(),
		aggregator: make(chan *Folder, 64),
	}
}

func (r *reader) Handler(path string) {
	for _, m := range r.modules {
		if m.allow(path) {
			m.read(path, r.aggregator)
		}
	}
}

func (r *reader) aggregate() {
	r.wg.Add(1)
	defer r.wg.Done()
	for folder := range r.aggregator {
		current, ok := r.result.Search(folder)
		if !ok {
			// Create new node
			r.result.Insert(folder)
		} else {
			// Update current, that has already been created before
			folder.copyContent(current.(*Folder))
		}
	}
}
