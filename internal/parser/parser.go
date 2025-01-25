package parser

type DirectoryParser interface {
	Parse(root string) ([]string, error)
	// FilterByExtenstion(files []string, ext string) []string
	GetFilesByTypes(allFiles []string, types []string) []string
}

type FileParser interface {
	ReadContent(path string) (string, error)
}
