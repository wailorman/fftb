package files

import "path/filepath"

// NewFile _
func NewFile(relativePath string) *File {
	fullPath, _ := filepath.Abs(relativePath)

	dirPath, fileName := filepath.Split(fullPath)

	return &File{
		fileName: fileName,
		dirPath:  dirPath,
	}
}

// NewPath _
func NewPath(relativePath string) *Path {
	fullPath, _ := filepath.Abs(relativePath)

	return &Path{
		path: fullPath,
	}
}
