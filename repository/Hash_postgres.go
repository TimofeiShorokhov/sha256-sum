//Actions with database
package repository

import (
	"database/sql"
	"log"
	"sha256-sum/models"
)

type HashPostgres struct {
	db *sql.DB
}

//Creating database object
func NewHashPostgres(db *sql.DB) *HashPostgres {
	return &HashPostgres{db: db}
}

func (r *HashPostgres) CheckDB() int {
	var code int

	count, err := r.db.Query("SELECT COUNT(*) FROM shasum;")

	if err != nil {
		log.Fatalf("error with db.Query, %s: ", err)
	}

	defer count.Close()

	var num int

	for count.Next() {
		if err := count.Scan(&num); err != nil {
			log.Println(err)
		}
	}

	if num == 0 {
		code = 0
	}
	if num > 0 {
		code = 1
	}
	return code
}

//Getting data of all hashes from database
func (r *HashPostgres) GetDataFromDB() ([]models.HashData, error) {
	var hashes []models.HashData

	selectValue := `Select file, checksum, file_path, algorithm from shasum;`

	get, err := r.db.Query(selectValue)

	if err != nil {
		log.Println("error of getting data: " + err.Error())
		return []models.HashData{}, err
	}

	for get.Next() {
		var hash models.HashData
		err = get.Scan(&hash.FileName, &hash.CheckSum, &hash.FilePath, &hash.Algorithm)
		hashes = append(hashes, hash)
	}
	return hashes, nil
}

//Inserting data in database
func (r *HashPostgres) PutDataInDB(data []models.HashData, podData models.PodData) error {

	transaction, err := r.db.Begin()

	if err != nil {
		log.Println("error with database: " + err.Error())
	}
	query := `INSERT INTO shasum(file,checksum,file_path,algorithm,pod_name,container_name,image_name,creation_time) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) 
ON CONFLICT ON CONSTRAINT shasum_unique DO UPDATE SET checksum=excluded.checksum`

	for _, h := range data {
		_, err := transaction.Exec(query, h.FileName, h.CheckSum, h.FilePath, h.Algorithm, podData.PodName, podData.ContainerName, podData.ImageName, podData.CreationTime)
		if err != nil {
			transaction.Rollback()
			return err
		}
	}
	return transaction.Commit()
}
func (r *HashPostgres) PutPodInDB(name string) error {

	transaction, err := r.db.Begin()

	if err != nil {
		log.Println("error with database: " + err.Error())
	}
	query := `INSERT INTO shasum(file,checksum,file_path,algorithm) VALUES ($1,$2,$3,$4) 
ON CONFLICT ON CONSTRAINT shasum_unique DO UPDATE SET checksum=excluded.checksum`
	_, err = transaction.Exec(query, name, "12345", "path", "sha256")
	if err != nil {
		transaction.Rollback()
		return err
	}
	return transaction.Commit()
}

func (r *HashPostgres) GetDataByPathFromDB(alg string) ([]models.HashData, error) {
	var hashes []models.HashData

	selectValue := `Select file, checksum, file_path, algorithm from shasum where algorithm = $1`

	get, err := r.db.Query(selectValue, alg)

	if err != nil {
		log.Println("error of getting data: " + err.Error())
		return []models.HashData{}, err
	}

	for get.Next() {
		var hash models.HashData
		err = get.Scan(&hash.FileName, &hash.CheckSum, &hash.FilePath, &hash.Algorithm)
		hashes = append(hashes, hash)
	}
	return hashes, nil
}

func (r *HashPostgres) UpdateDeletedStatusInDB(data []models.HashData) error {

	transaction, err := r.db.Begin()

	if err != nil {
		log.Println("error with database: " + err.Error())
	}
	query := `UPDATE shasum SET deleted = true where file_path = $1 and algorithm = $2`

	for _, h := range data {
		_, err := transaction.Exec(query, h.FilePath, h.Algorithm)
		if err != nil {
			transaction.Rollback()
			return err
		}
	}
	return transaction.Commit()
}

func (r *HashPostgres) Truncate() error {

	_, err := r.db.Exec("TRUNCATE TABLE shasum;")

	if err != nil {
		log.Fatalln(err)
	}
	return nil
}
