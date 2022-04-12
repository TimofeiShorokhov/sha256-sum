package main

import (
	"flag"
	"fmt"
	"log"
	"sha256-sum/repository"
	"sha256-sum/services"
	"time"
)

var (
	filePath string
	helpPath bool
	dirPath  string
	getData  bool
	hashAlg  string
)

func init() {
	flag.StringVar(&filePath, "p", "", "Hash through file path")
	flag.StringVar(&dirPath, "d", "", "Hash of all files through directory path")
	flag.StringVar(&hashAlg, "a", "", "md5, sha512, default: sha256")
	flag.BoolVar(&getData, "g", false, "Get all data from database")
	flag.BoolVar(&helpPath, "h", false, "info")
	flag.Parse()
}

func main() {
	start := time.Now()
	database, err := repository.NewPostgresDB(repository.PostgresDB{
		Host:     "localhost",
		Port:     "5432",
		User:     "tim",
		Password: "123",
		DBName:   "timdb",
		SSLMode:  "disable"})
	if err != nil {
		log.Fatal("failed to initialize dao:", err.Error())
	}
	repository := repository.NewRepository(database)
	ser := services.NewService(repository)

	services.CatchStopSignal()
	ser.CallFunction(filePath, helpPath, dirPath, getData, hashAlg)
	fmt.Println(time.Since(start).Seconds())
}
