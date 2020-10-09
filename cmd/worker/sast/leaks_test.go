package sast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObfuscate(t *testing.T) {
	str := "testing_obfuscate"
	expected := "testing_obf******"
	actual := obfuscate(str)

	assert.Equal(t, expected, actual)
}
