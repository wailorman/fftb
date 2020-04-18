package ffchunker

import "github.com/wailorman/ffchunker/files"

// Chunker _
type Chunker struct {
	mainFile   *files.Filer
	resultPath string
}

// Start _
func (c *Chunker) Start() error {
	return nil
}
