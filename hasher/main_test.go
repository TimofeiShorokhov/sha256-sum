package hasher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUnknown(t *testing.T) {
	hasher, err := New("sha1024")
	if assert.Nil(t, hasher) && assert.Error(t, err) {
		assert.Equal(t, "no such algorithm", err.Error())
	}
}
