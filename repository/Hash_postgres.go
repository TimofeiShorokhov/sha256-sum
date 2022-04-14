//Actions with database
package repository

import (
	"database/sql"
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

//Creating database object
func NewHashPostgres(db *sql.DB) *HashPostgres {
	return &HashPostgres{db: db}
}

//Getting data of all hashes from database
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

//Inserting data in database
func (r *HashPostgres) PutDataInDB(data []HashData) error {

	transaction, err := r.db.Begin()

	if err != nil {
		log.Println("error with database: " + err.Error())
	}
	query := `INSERT INTO shasum(file,checksum,file_path,algorithm) VALUES ($1,$2,$3,$4) 
ON CONFLICT ON CONSTRAINT shasum_unique DO UPDATE SET checksum=excluded.checksum`

	for _, h := range data {
		_, err := transaction.Exec(query, h.FileName, h.CheckSum, h.FilePath, h.Algorithm)
		if err != nil {
			transaction.Rollback()
		}
	}
	return transaction.Commit()
}
