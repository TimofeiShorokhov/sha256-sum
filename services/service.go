package services

import (
	"sha256-sum/repository"
	"sync"
)

type HashApp interface {
	GetData() ([]repository.HashData, error)
	PutData(res HashDataUtils) (int, error)
	Worker(wg *sync.WaitGroup, jobs <-chan string, results chan<- HashDataUtils, hashAlg string)
	CheckSum(path string, hashAlg string)
	CallFunction(filePath string, helpPath bool, dirPath string, getData bool, hashAlg string)
}

type Service struct {
	HashApp
}

func NewService(rep *repository.Repository) *Service {
	return &Service{
		NewHashService(*rep),
	}
}
