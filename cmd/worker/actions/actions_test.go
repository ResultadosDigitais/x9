package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRepoInfo(t *testing.T) {
	repository := "https://github.com/ResultadosDigitais/x9.git"
	expectedOwner := "ResultadosDigitais"
	expectedRepo := "x9"
	owner, repo := getRepoInfo(repository)
	assert.Equal(t, expectedOwner, owner)
	assert.Equal(t, expectedRepo, repo)
}
