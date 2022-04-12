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
	PutDataInDB(fileName string, checksum string, filePath string, algorithm string) (int, error)
}
