package tests

import (
	"github.com/magiconair/properties/assert"
	assert2 "github.com/stretchr/testify/assert"
	"sha256-sum/services"
	"testing"
)

func TestHasher(t *testing.T) {

	//hash of file by md5 algorithm
	hash := "0c5856fc21e958f96f79b95d279ad60a"

	//random string for testing not equal func
	notEqualHash := "0c5856fcs3e558f96f79b95d279ad60a"

	//hash from func by md5 algorithm
	hashFromFunc := services.HashOfFile("/home/tshorokhov@scnsoft.com/Pictures/1", "md5")

	//testing HashOfFile function for equal
	assert.Equal(t, hash, hashFromFunc.Checksum)

	//testing HashOfFile function for not equal
	assert2.NotEqual(t, notEqualHash, hashFromFunc.Checksum)

}

/*
func TestDir(t *testing.T) {
	paths := make(chan string)
	pathsForFunc := make(chan string)
	var pathes []string
	pathes = append(pathes, "/home/tshorokhov@scnsoft.com/Pictures/pictures/Screenshot from 2022-03-30 14-51-18.png",
		"/home/tshorokhov@scnsoft.com/Pictures/pictures/Screenshot from 2022-03-30 14-51-34.png",
		"/home/tshorokhov@scnsoft.com/Pictures/pictures/Screenshot from 2022-03-31 12-39-34.png",
		"/home/tshorokhov@scnsoft.com/Pictures/pictures/Screenshot from 2022-03-31 12-40-05.png",
		"/home/tshorokhov@scnsoft.com/Pictures/pictures/Screenshot from 2022-03-31 14-42-22.png")
	go func() {
		for _, h := range pathes {
			paths <- h
		}
		close(paths)
	}()
	go services.HashOfDir("/home/tshorokhov@scnsoft.com/Pictures/pictures", pathsForFunc)
	var pa []string
	for i := 0; i < 5; i++ {
		pa = append(pa, <-pathsForFunc)
	}
	for _, h := range pathes {
		for _, h1 := range pa {
			assert.Equal(t, h1, h)
		}
	}
}

*/
