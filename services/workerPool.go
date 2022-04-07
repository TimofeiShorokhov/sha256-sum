package services

import (
	"context"
	"fmt"
	"os"
	"sync"
)

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
