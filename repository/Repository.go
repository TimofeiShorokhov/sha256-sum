package repository

import (
	"database/sql"
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
	GetDataFromDB() ([]HashData, error)
	PutDataInDB(data []HashData) error
	GetDataByPathFromDB(dir string, alg string) ([]HashData, error)
}
