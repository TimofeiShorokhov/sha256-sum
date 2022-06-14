package hasher

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNewMD5(t *testing.T) {
	assert.IsType(t, &MD5{}, NewMD5())
}

func TestHasher(t *testing.T) {

	hashService := NewMD5()

	r := strings.NewReader("dsfdsa")
	hash, err := hashService.Hash(r)

	if assert.Nil(t, err) {
		assert.Equal(t, "18f1b29e9a03f3b8bc2cd8233057b28a", hash)
	}
}
