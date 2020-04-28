package sast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsBlacklisted(t *testing.T) {

	blackListedRepositories := []string{
		"OrganizationAccount/RepositoryName2",
		"OrganizationAccount/RepositoryName3",
		"AnotherOrganizationAccount/RepositoryName",
	}
	repo := "http://github.com/OrganizationAccount/RepositoryName2/"

	current := isBlacklisted(repo, blackListedRepositories)
	assert.True(t, current)

	repo = "http://github.com/AnotherOrganizationAccount/RepositoryName2/"

	current = isBlacklisted(repo, blackListedRepositories)
	assert.False(t, current)

	repo = "http://github.com/OrganizationAccount/RepositoryName/"

	current = isBlacklisted(repo, blackListedRepositories)
	assert.False(t, current)
}

func TestGetHash(t *testing.T) {
	actual := getHash("abc")
	expected := "a9993e364706816aba3e25717850c26c9cd0d89d"

	assert.Equal(t, expected, actual)

	actual = getHash("abcdefgh", "18ncyas8d", "a", "(*q/qwe")
	expected = "5a74a7163108aca1983083127fbc3f9e65ce3ffe"

	assert.Equal(t, expected, actual)
}
