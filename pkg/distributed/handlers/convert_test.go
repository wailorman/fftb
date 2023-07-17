package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test__searchOutputDirs(tg *testing.T) {
	tg.Run("opts with output dirs", func(t *testing.T) {
		actual := searchOutputDirs("/home", []string{
			"-i", "abc",
			"-preset",
			"[abc]",
			"output/WO/example.mp4",
			"output/WO/abc/example2.mp4",
		})

		assert.Equal(t, []string{"/home/output/WO", "/home/output/WO/abc"}, actual)
	})

	tg.Run("opts without output dirs", func(t *testing.T) {
		actual := searchOutputDirs("/home", []string{
			"-i", "abc",
			"-preset",
			"[abc]",
		})

		assert.Equal(t, []string{}, actual)
	})
}
