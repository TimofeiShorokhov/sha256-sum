package services

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"sha256-sum/errors"
)

func HashOfFile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		errors.CheckErr(err)
		return ""
	}

	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		errors.CheckErr(err)
		return ""
	}
	res := "File: " + file.Name() + ", Checksum: " + hex.EncodeToString(hash.Sum(nil))
	return res
}
