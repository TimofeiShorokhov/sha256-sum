package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
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

/*
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
	ser := services.NewService(repository, hashAlg)

	services.CatchStopSignal()
	ser.CallFunction(helpPath, dirPath, getData, getChangedData, updDeleted)
	fmt.Println(time.Since(start).Seconds())
}


*/

func main() {
	start := time.Now()

	hashAlg = "sha256"
	path := "/home/tshorokhov@scnsoft.com/Pictures"

	cfg, err := configs.ParseConfig()

	if err != nil {
		log.Fatal("error parsing config: " + err.Error())
	}

	database, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("failed to initialize dao:", err.Error())
	}
	repository := repository.NewRepository(database)
	ser := services.NewService(repository, hashAlg)

	c := make(chan os.Signal, 1)
	signal.Notify(c)

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				ser.Operations(repository.CheckDB(), path)
			}
		}
	}()

	<-c
	ticker.Stop()

	services.CatchStopSignal()
	fmt.Println(time.Since(start).Seconds())
}
