package services

import (
	"context"
	"sha256-sum/repository"
	"sync"
)

type HashApp interface {
	GetData() ([]repository.HashData, error)
	GetChangedData(dir string) error
	PutData(res []HashDataUtils) error
	Worker(wg *sync.WaitGroup, jobs <-chan string, results chan<- HashDataUtils)
	CheckSum(path string) []HashDataUtils
	CallFunction(helpPath bool, dirPath string, getData bool, getChangedData string, updDeleted string)
	Result(ctx context.Context, results chan HashDataUtils) []HashDataUtils
	UpdateDeletedStatus(dir string) error
	HashOfFile(path string) HashDataUtils
}

type Service struct {
	HashApp
}

func NewService(rep *repository.Repository, algo string) *Service {
	return &Service{
		NewHashService(*rep, algo),
	}
}
