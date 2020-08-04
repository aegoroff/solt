package msvc

// ReaderHandler defines file system scanning handler
type ReaderHandler interface {
	// Handler method called on each file and folder scanned
	Handler(path string)
}

type readerModule interface {
	filter(path string) bool
	read(path string) (*Folder, bool)
}

type walker interface {
	onFolder(f *Folder)
}
