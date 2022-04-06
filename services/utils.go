package services

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sha256-sum/errors"
	"sync"
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
	res := fmt.Sprintf("%x %s", hash.Sum(nil), filepath.Base(path))
	return res
}

func HashOfDir(path string, paths chan string) {

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				errors.CheckErr(err)
				return err
			}
			if info.IsDir() == false {
				paths <- path
			}
			return nil
		})
	if err != nil {
		errors.CheckErr(err)
	}
	close(paths)
}

func Worker(wg *sync.WaitGroup, jobs <-chan string, results chan<- string) {
	defer wg.Done()
	for j := range jobs {
		results <- HashOfFile(j)
	}
}

func Result(results chan string) {

	for {
		select {
		case hash, ok := <-results:
			if !ok {
				return
			}
			fmt.Println(hash)
		}
	}
}

func CheckSum(path string) {
	jobs := make(chan string)
	results := make(chan string)

	go HashOfDir(path, jobs)
	go func() {
		var wg sync.WaitGroup
		for w := 1; w <= 10; w++ {
			wg.Add(1)
			go Worker(&wg, jobs, results)
		}
		defer close(results)
		wg.Wait()
	}()
	Result(results)
}
