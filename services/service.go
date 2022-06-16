package services

import (
	"sha256-sum/repository"
	"sync"
)

type HashApp interface {
	GetData() ([]repository.HashData, error)
	GetChangedData(dir string) (int, error)
	PutData(res []HashDataUtils) error
	Worker(wg *sync.WaitGroup, jobs <-chan string, results chan<- HashDataUtils)
	CheckSum(path string) []HashDataUtils
	CallFunction(helpPath bool, dirPath string, getData bool, getChangedData string, updDeleted string)
	Result(results chan HashDataUtils) []HashDataUtils
	UpdateDeletedStatus(dir string) error
	HashOfFile(path string) HashDataUtils
	SavingData(data []HashDataUtils)
	Operations(code int, path string)
}

type Service struct {
	HashApp
}

func NewService(rep *repository.Repository, algo string) *Service {
	return &Service{
		NewHashService(*rep, algo),
	}
}
