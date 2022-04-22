package main

import (
	"flag"
	"fmt"
	"log"
	"sha256-sum/configs"
	"sha256-sum/repository"
	"sha256-sum/services"
	"time"
)

var (
	helpPath       bool
	dirPath        string
	getData        bool
	getChangedData string
	updDeleted     string
	hashAlg        string
)

//Creating and parsing flags
func init() {
	flag.StringVar(&dirPath, "d", "", "Hash of all files through directory path")
	flag.StringVar(&hashAlg, "a", "", "md5, sha512, default: sha256")
	flag.StringVar(&getChangedData, "c", "", "Check of changed checksum")
	flag.StringVar(&updDeleted, "u", "", "Update deleted status in database")
	flag.BoolVar(&getData, "g", false, "Get all data from database")
	flag.BoolVar(&helpPath, "h", false, "info")
	flag.Parse()
}

func main() {
	start := time.Now()

	cfg, err := configs.ParseConfig()

	if err != nil {
		log.Fatal("error parsing config: " + err.Error())
	}

	database, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("failed to initialize dao:", err.Error())
	}
	repository := repository.NewRepository(database)
	ser := services.NewService(repository)

	services.CatchStopSignal()
	ser.CallFunction(helpPath, dirPath, getData, getChangedData, updDeleted, hashAlg)
	fmt.Println(time.Since(start).Seconds())
}
