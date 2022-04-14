package services

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sha256-sum/errors"
	"sha256-sum/repository"
)

type HashService struct {
	repo repository.Repository
}

func NewHashService(repo repository.Repository) *HashService {
	return &HashService{
		repo: repo,
	}
}

type HashDataUtils struct {
	FileName  string
	Checksum  string
	FilePath  string
	Algorithm string
}

//Generating checksum of file
func HashOfFile(path string, hashAlg string) HashDataUtils {
	var res HashDataUtils

	file, err := os.Open(path)
	if err != nil {
		errors.CheckErr(err)
		return HashDataUtils{}
	}

	defer file.Close()

	var checkSum string

	switch hashAlg {
	case "md5":
		hash := md5.New()
		_, err = io.Copy(hash, file)
		checkSum = hex.EncodeToString(hash.Sum(nil))
	case "sha512":
		hash := sha512.New()
		_, err = io.Copy(hash, file)
		checkSum = hex.EncodeToString(hash.Sum(nil))
	default:
		hash := sha256.New()
		hashAlg = "sha256"
		_, err = io.Copy(hash, file)
		checkSum = hex.EncodeToString(hash.Sum(nil))
	}
	if err != nil {
		errors.CheckErr(err)
		return HashDataUtils{}
	}
	res.FileName = filepath.Base(path)
	res.Checksum = checkSum
	res.FilePath = file.Name()
	res.Algorithm = hashAlg

	return res
}

//Generating checksum of all files by the directory
func HashOfDir(path string, paths chan string) {
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				errors.CheckErr(err)
				return err
			}
			if info.IsDir() == false {
				paths <- path
			}
			return nil
		})
	if err != nil {
		errors.CheckErr(err)
	}
	close(paths)
}

//Custom interrupt
func CatchStopSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("stopped by user %d", sig)
			os.Exit(1)
		}
	}()
}

//Function for calling checksum function
func (s *HashService) CallFunction(filePath string, helpPath bool, dirPath string, getData bool, getChangedData bool, hashAlg string) {
	switch {
	case helpPath:
		flag.Usage = func() {
			fmt.Fprintf(os.Stderr, "Help with commands %s:\nUse one of the following commands:\n", os.Args[0])
			flag.VisitAll(func(f *flag.Flag) {
				fmt.Fprintf(os.Stderr, " flag:	-%v \n 	%v\n", f.Name, f.Usage)
			})
		}
		flag.Usage()
	case len(filePath) > 0:
		fmt.Println(HashOfFile(filePath, hashAlg))
	case len(dirPath) > 0:
		s.SavingData(s.CheckSum(dirPath, hashAlg))
	case getData:
		s.GetData()
	default:
		log.Println("Error with flag, use '-h' flag for help ")
	}
}

//Getting data
func (s *HashService) GetData() ([]repository.HashData, error) {
	data, err := s.repo.GetDataFromDB()

	if data == nil {
		log.Println("no data for output")
		return nil, err
	}

	if err != nil {
		fmt.Println(err)
	}
	for _, h := range data {
		fmt.Printf("File name: %s, Checksum: %s, Algorithm: %s\n", h.FileName, h.CheckSum, h.Algorithm)
	}
	return data, nil
}

//Inserting data
func (s *HashService) PutData(res []HashDataUtils) error {
	var data []repository.HashData
	var dat repository.HashData
	for _, h := range res {
		dat.FileName = h.FileName
		dat.FilePath = h.FilePath
		dat.Algorithm = h.Algorithm
		dat.CheckSum = h.Checksum
		data = append(data, dat)
	}
	return s.repo.PutDataInDB(data)
}
