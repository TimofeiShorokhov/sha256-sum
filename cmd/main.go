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
	filePath       string
	helpPath       bool
	dirPath        string
	getData        bool
	getChangedData bool
	hashAlg        string
)

func init() {
	flag.StringVar(&filePath, "p", "", "Hash through file path")
	flag.StringVar(&dirPath, "d", "", "Hash of all files through directory path")
	flag.StringVar(&hashAlg, "a", "", "md5, sha512, default: sha256")
	flag.BoolVar(&getData, "g", false, "Get all data from database")
	flag.BoolVar(&getChangedData, "c", false, "Get all changed data from database")
	flag.BoolVar(&helpPath, "h", false, "info")
	flag.Parse()
}

//Hello there
func main() {
	start := time.Now()

	cfg, err := configs.ParseConfig("configs/")

	if err != nil {
		log.Println("error: " + err.Error())
	}

	database, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("failed to initialize dao:", err.Error())
	}
	repository := repository.NewRepository(database)
	ser := services.NewService(repository)

	services.CatchStopSignal()
	ser.CallFunction(filePath, helpPath, dirPath, getData, getChangedData, hashAlg)
	fmt.Println(time.Since(start).Seconds())
}
