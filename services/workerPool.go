package services

import (
	"fmt"
	"sha256-sum/models"
	"sync"
)

//Workers for goroutines
func (s *HashService) Worker(wg *sync.WaitGroup, jobs <-chan string, results chan<- models.HashDataUtils) {
	defer wg.Done()
	for j := range jobs {
		results <- s.HashOfFile(j)
	}
}

//Inserting results of hashing in slice of structures and printing
func (s *HashService) Result(results chan models.HashDataUtils) []models.HashDataUtils {
	var data []models.HashDataUtils
	for {
		select {
		case hash, ok := <-results:
			if !ok {
				return data
			}
			data = append(data, hash)
		}

	}
}

//Inserting data
func (s *HashService) SavingData(data []models.HashDataUtils, podData models.PodInfo) {
	for _, h := range data {
		fmt.Printf("File name: %s, Checksum: %s, Algorithm: %s\n", h.FileName, h.Checksum, h.Algorithm)
	}
	s.PutData(data, podData)
}

//Hashing files with workers
func (s *HashService) CheckSum(path string) []models.HashDataUtils {

	jobs := make(chan string)
	results := make(chan models.HashDataUtils)

	go HashOfDir(path, jobs)
	go func() {
		var wg sync.WaitGroup
		for w := 1; w <= 10; w++ {
			wg.Add(1)
			go s.Worker(&wg, jobs, results)
		}
		defer close(results)
		wg.Wait()
	}()
	return s.Result(results)
}
