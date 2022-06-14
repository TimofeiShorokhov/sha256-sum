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
	CheckDB() int
	GetDataFromDB() ([]HashData, error)
	PutDataInDB(data []HashData) error
	GetDataByPathFromDB(alg string) ([]HashData, error)
	UpdateDeletedStatusInDB(data []HashData) error
	Truncate() error
}
