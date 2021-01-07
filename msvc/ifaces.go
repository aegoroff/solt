package msvc

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

// Container provides paths container interface
type Container interface {
	// Path provides Container's path itself
	Path() string
	// Items provides all paths included into container
	Items() []string
}
