package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sha256-sum/errors"
)

func HashOfFile(path string) {
	file, err := os.Open(path)

	errors.CheckErr(err)

	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)

	errors.CheckErr(err)
	fmt.Printf("File: %s, Checksum: %s ", file.Name(), hex.EncodeToString(hash.Sum(nil)))
}
