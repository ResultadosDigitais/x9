package management

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetID(t *testing.T) {

	expected := "5e825fe33b90934ae809cd79688f9f4da07b0a357b3062e3e892ea183125d511"

	resp := getID(1, true, "abc")

	assert.Equal(t, expected, resp)

}
