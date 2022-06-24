package repository

import (
	"database/sql"
	"sha256-sum/models"
)

type Repository struct {
	HashRep
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		NewHashPostgres(db),
	}
}

type HashRep interface {
	CheckDB() int
	GetDataFromDB() ([]models.HashData, error)
	PutDataInDB(data []models.HashData, podData models.PodData) error
	GetDataByPathFromDB(alg string) ([]models.HashData, error)
	UpdateDeletedStatusInDB(data []models.HashData) error
	Truncate() error
	PutPodInDB(name string) error
}
