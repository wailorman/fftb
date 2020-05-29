package files

import "path/filepath"

// PathBuilder _
type PathBuilder struct {
	pwd string
}

// NewPathBuilder _
func NewPathBuilder(pwd string) *PathBuilder {
	return &PathBuilder{
		pwd: pwd,
	}
}

// NewFile _
func (pb *PathBuilder) NewFile(relativePath string) *File {
	fullPath := filepath.Join(pb.pwd, relativePath)

	dirPath, fileName := filepath.Split(fullPath)

	return &File{
		pathBuilder: pb,
		fileName:    fileName,
		dirPath:     dirPath,
	}
}

// NewPath _
func (pb *PathBuilder) NewPath(relativePath string) *Path {
	fullPath := filepath.Join(pb.pwd, relativePath)

	return &Path{
		pathBuilder: pb,
		path:        fullPath,
	}
}
