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
	transaction, err := r.db.Begin()

	if err != nil {
		log.Println("error with database: " + err.Error())
	}

	defer transaction.Commit()

	insertValue := `INSERT INTO shasum(file, checksum, file_path, algorithm) VALUES ($1,$2,$3,$4)`

	row := transaction.QueryRow(insertValue, fileName, checksum, filePath, algorithm)

	if err = row.Scan(&HashId); err != nil {
		return 0, fmt.Errorf("error while scanning for id: %s", err)
	}
	return HashId, nil
}

func (r *HashPostgres) GetChangedHashFromDB() {

	selectValus := `SELECT file,algorithm,count(*) FROM shasum group by file,algorithm having count(*) > 1`

	get, err := r.db.Query(selectValus)
	if err != nil {
		log.Println("error of getting data: " + err.Error())
	}

	defer get.Close()

	for get.Next() {
		var file string
		var alg string
		var count int
		err = get.Scan(&file, &alg, &count)
		fmt.Printf("Checksum of this file: %s, with algorithm: %s, was changed\n", file, alg)
	}
}
