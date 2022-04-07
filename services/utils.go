package services

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sha256-sum/errors"
	"sync"
)

func HashOfFile(path string, hashAlg string) string {

	file, err := os.Open(path)
	if err != nil {
		errors.CheckErr(err)
		return ""
	}

	defer file.Close()

	var checkSum interface{}

	switch hashAlg {
	case "md5":
		hash := md5.New()
		_, err = io.Copy(hash, file)
		checkSum = hash.Sum(nil)
	case "512":
		hash := sha512.New()
		_, err = io.Copy(hash, file)
		checkSum = hash.Sum(nil)
	default:
		hash := sha256.New()
		_, err = io.Copy(hash, file)
		checkSum = hash.Sum(nil)
	}
	if err != nil {
		errors.CheckErr(err)
		return ""
	}
	res := fmt.Sprintf("%x %s", checkSum, filepath.Base(path))
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

func Worker(wg *sync.WaitGroup, jobs <-chan string, results chan<- string, hashAlg string) {
	defer wg.Done()
	for j := range jobs {
		results <- HashOfFile(j, hashAlg)
	}
}

func Result(ctx context.Context, results chan string) {

	for {
		select {
		case hash, ok := <-results:
			if !ok {
				return
			}
			fmt.Println(hash)
		case <-ctx.Done():
			fmt.Println("canceled by user")
			os.Exit(1)
			return
		}
	}
}

func CheckSum(path string, hashAlg string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		fmt.Scanln()
		cancel()
	}()

	jobs := make(chan string)
	results := make(chan string)

	go HashOfDir(path, jobs)
	go func() {
		var wg sync.WaitGroup
		for w := 1; w <= 10; w++ {
			wg.Add(1)
			go Worker(&wg, jobs, results, hashAlg)
		}
		defer close(results)
		wg.Wait()
	}()
	Result(ctx, results)
}
