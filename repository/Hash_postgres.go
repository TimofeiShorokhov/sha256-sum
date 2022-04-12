package repository

import (
	"database/sql"
	"fmt"
	"log"
)

type HashPostgres struct {
	db *sql.DB
}

type HashData struct {
	Id        int
	FileName  string
	CheckSum  string
	FilePath  string
	Algorithm string
}

func NewHashPostgres(db *sql.DB) *HashPostgres {
	return &HashPostgres{db: db}
}

func (r *HashPostgres) GetDataFromDB() ([]HashData, error) {
	var hashes []HashData

	selectValue := `Select file, checksum, file_path, algorithm from shasum`

	get, err := r.db.Query(selectValue)

	if err != nil {
		log.Println("error of getting data: " + err.Error())
		return []HashData{}, err
	}

	for get.Next() {
		var hash HashData
		err = get.Scan(&hash.FileName, &hash.CheckSum, &hash.FilePath, &hash.Algorithm)
		hashes = append(hashes, hash)
	}
	return hashes, nil
}

func (r *HashPostgres) PutDataInDB(fileName string, checksum string, filePath string, algorithm string) (int, error) {
	var HashId int
	insertValue := `INSERT INTO shasum(file, checksum, file_path, algorithm) VALUES ($1,$2,$3,$4)
	   on conflict (file) do update set checksum=excluded.checksum,algorithm=excluded.algorithm `

	row := r.db.QueryRow(insertValue, fileName, checksum, filePath, algorithm)

	if err := row.Scan(&HashId); err != nil {
		return 0, fmt.Errorf("error while scanning for id: %s", err)
	}
	return HashId, nil
}