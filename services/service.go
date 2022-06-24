package services

import (
	"k8s.io/client-go/kubernetes"
	"sha256-sum/models"
	"sha256-sum/repository"
	"sync"
)

type HashApp interface {
	GetData() ([]models.HashData, error)
	GetChangedData(dir string) (int, error)
	PutData(res []models.HashDataUtils, podData models.PodInfo) error
	Worker(wg *sync.WaitGroup, jobs <-chan string, results chan<- models.HashDataUtils)
	CheckSum(path string) []models.HashDataUtils
	CallFunction(helpPath bool, dirPath string, getData bool, getChangedData string, updDeleted string, podData models.PodInfo)
	Result(results chan models.HashDataUtils) []models.HashDataUtils
	UpdateDeletedStatus(dir string) error
	HashOfFile(path string) models.HashDataUtils
	SavingData(data []models.HashDataUtils, podData models.PodInfo)
	Operations(code int, path string)
	Podkicker(code int, path string)
	ConnToPod() *kubernetes.Clientset
}

type Service struct {
	HashApp
}

func NewService(rep *repository.Repository, algo string) *Service {
	return &Service{
		NewHashService(*rep, algo),
	}
}
