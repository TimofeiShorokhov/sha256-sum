package services

import (
	"context"
	"fmt"
	"os"
	"sync"
)

//Workers for goroutines
func (s *HashService) Worker(wg *sync.WaitGroup, jobs <-chan string, results chan<- HashDataUtils, hashAlg string) {
	defer wg.Done()
	for j := range jobs {
		results <- HashOfFile(j, hashAlg)
	}
}

//Inserting results of hashing in slice of structures and printing
func (s *HashService) Result(ctx context.Context, results chan HashDataUtils) []HashDataUtils {
	var data []HashDataUtils
	for {
		select {
		case hash, ok := <-results:
			if !ok {
				return data
			}
			data = append(data, hash)

		case <-ctx.Done():
			fmt.Println("canceled by user")
			os.Exit(1)
			return []HashDataUtils{}
		}

	}
}

//Inserting data
func (s *HashService) SavingData(data []HashDataUtils) {
	for _, h := range data {
		fmt.Printf("File name: %s, Checksum: %s, Algorithm: %s\n", h.FileName, h.Checksum, h.Algorithm)
	}
	s.PutData(data)
}

//Hashing files with workers
func (s *HashService) CheckSum(path string, hashAlg string) []HashDataUtils {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		fmt.Scanln()
		cancel()
	}()

	jobs := make(chan string)
	results := make(chan HashDataUtils)

	go HashOfDir(path, jobs)
	go func() {
		var wg sync.WaitGroup
		for w := 1; w <= 10; w++ {
			wg.Add(1)
			go s.Worker(&wg, jobs, results, hashAlg)
		}
		defer close(results)
		wg.Wait()
	}()
	return s.Result(ctx, results)
}
