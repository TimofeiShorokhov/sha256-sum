package services

import (
	"github.com/magiconair/properties/assert"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestHasher(t *testing.T) {

	//hash of file by md5 algorithm
	hash := "0c5856fc21e958f96f79b95d279ad60a"

	//random string for testing not equal func
	notEqualHash := "0c5856fcs3e558f96f79b95d279ad60a"

	//hash from func by md5 algorithm
	hashFromFunc := HashOfFile("/home/tshorokhov@scnsoft.com/Pictures/1", "md5")

	//testing HashOfFile function for equal
	assert.Equal(t, hash, hashFromFunc.Checksum)

	//testing HashOfFile function for not equal
	assert2.NotEqual(t, notEqualHash, hashFromFunc.Checksum)

}
