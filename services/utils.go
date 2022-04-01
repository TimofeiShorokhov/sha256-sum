package services

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
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

func HashOfDir(path string) []string {
	var res []string
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() == false {
				res = append(res, HashOfFile(path))
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return res
}
