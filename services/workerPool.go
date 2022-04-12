package services

import (
	"context"
	"fmt"
	"os"
	"sync"
)

func (s *HashService) Worker(wg *sync.WaitGroup, jobs <-chan string, results chan<- HashDataUtils, hashAlg string) {
	defer wg.Done()
	for j := range jobs {
		results <- HashOfFile(j, hashAlg)
		res := HashOfFile(j, hashAlg)
		s.PutData(res)
	}
}

func (s *HashService) Result(ctx context.Context, results chan HashDataUtils) {

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

func (s *HashService) CheckSum(path string, hashAlg string) {
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
	s.Result(ctx, results)
}
