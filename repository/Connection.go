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

func NewPostgresDB(cfg *PostgresDB) (*sql.DB, error) {
	pgsqlConn := fmt.Sprintf("host= %s port= %s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)
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
