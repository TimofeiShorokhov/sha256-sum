package main

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sha256-sum/repository"
	"sha256-sum/services"
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
	hashAlg = "sha256"
	var err error
	err = godotenv.Load()

	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	}
	database, err := repository.ConnToDb(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	if err != nil {
		log.Fatal("failed to initialize dao:", err.Error())
	}

	repository := repository.NewRepository(database)
	ser := services.NewService(repository, hashAlg)
	ser.Operations(repository.CheckDB(), "/root")
}
