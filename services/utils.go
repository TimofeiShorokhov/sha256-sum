package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sha256-sum/errors"
	"sync"
)

const help = "\nUse -p flag to calculate checksum for file \nUse -d flag to calculate checksum of all files in directory"

func HashOfFile(path string, wg *sync.WaitGroup) string {
	defer wg.Done()

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

func HashOfDir(path string, wg *sync.WaitGroup) {
	defer wg.Done()

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				errors.CheckErr(err)
				return err
			}
			if info.IsDir() == false {
				wg.Add(1)
				go fmt.Println(HashOfFile(path, wg))
			}
			return nil
		})
	if err != nil {
		errors.CheckErr(err)
	}
}

func PrintHelp(path string, wg *sync.WaitGroup) {
	defer wg.Done()
	if path == "--help" {
		log.Println(help)
	}
}
