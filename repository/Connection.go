package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type PostgresDB struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(dbs PostgresDB) (*sql.DB, error) {
	pgsqlConn := fmt.Sprintf("host= %s port= %s user=%s password=%s dbname=%s sslmode=disable", dbs.Host, dbs.Port, dbs.User, dbs.Password, dbs.DBName)
	db, err := sql.Open("postgres", pgsqlConn)
	if err != nil {

		return nil, fmt.Errorf("error connecting to database:%s", err)
	}
	err = db.Ping()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return db, nil
}
