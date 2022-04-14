package services

import (
	"context"
	"sha256-sum/repository"
	"sync"
)

type HashApp interface {
	GetData() ([]repository.HashData, error)
	PutData(res []HashDataUtils) error
	Worker(wg *sync.WaitGroup, jobs <-chan string, results chan<- HashDataUtils, hashAlg string)
	CheckSum(path string, hashAlg string) []HashDataUtils
	CallFunction(filePath string, helpPath bool, dirPath string, getData bool, getChangedData bool, hashAlg string)
	Result(ctx context.Context, results chan HashDataUtils) []HashDataUtils
}

type Service struct {
	HashApp
}

func NewService(rep *repository.Repository) *Service {
	return &Service{
		NewHashService(*rep),
	}
}
