package sast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getConfig() Config {
	return Config{
		BlacklistedExtensions:        []string{".ext"},
		BlacklistedPaths:             []string{"path/anotherpath"},
		BlacklistedEntropyExtensions: []string{".abc", ".def"},
	}
}

func TestIsSkippableFile(t *testing.T) {
	config := getConfig()
	path := "abc/path/anotherpath/file.abc"
	assert.True(t, IsSkippableFile(config, path))

	path = "abc/def/ghi.ext"
	assert.True(t, IsSkippableFile(config, path))

	path = "abc/def.extension"
	assert.False(t, IsSkippableFile(config, path))
}

func TestCanCheckEntropy(t *testing.T) {
	config := getConfig()

	match := MatchFile{
		Path:      "/xyz",
		Filename:  "jkl",
		Extension: ".mno",
	}

	assert.True(t, match.CanCheckEntropy(config))

	match.Extension = ".abc"
	assert.False(t, match.CanCheckEntropy(config))

}
