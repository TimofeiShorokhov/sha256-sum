package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type PostgresDB struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

//Connection to database
func NewPostgresDB(cfg *PostgresDB) (*sql.DB, error) {
	pgsqlConn := fmt.Sprintf("host= %s port= %s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := sql.Open("postgres", pgsqlConn)
	if err != nil {

		return nil, fmt.Errorf("error connecting to database:%s", err)
	}
	return db, nil
}

func ConnToDb(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) (*sql.DB, error) {
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	db, err := sql.Open(Dbdriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", Dbdriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database\n", Dbdriver)
	}
	return db, nil
}
