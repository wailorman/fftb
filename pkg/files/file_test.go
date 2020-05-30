package files

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test__NewWithSuffix(t *testing.T) {
	assert := assert.New(t)

	origFile := NewFile("filename.txt")
	suffixedFile := origFile.NewWithSuffix("_suffix")

	assert.Equal(origFile.Name(), "filename.txt")
	assert.Equal(suffixedFile.Name(), "filename_suffix.txt")
}
