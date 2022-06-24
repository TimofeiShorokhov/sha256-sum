package services

import (
	"flag"
	"fmt"
	"sha256-sum/models"

	//"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sha256-sum/errors"
	"sha256-sum/hasher"
	"sha256-sum/repository"
)

type HashService struct {
	repo   repository.Repository
	hasher hasher.Hasher
	alg    string
}

func NewHashService(repo repository.Repository, algo string) *HashService {
	h, err := hasher.New(algo)
	if err != nil {
		log.Fatalf("err: %s", err)
	}
	return &HashService{
		repo:   repo,
		hasher: h,
		alg:    algo,
	}
}

//Generating checksum of file
func (s *HashService) HashOfFile(path string) models.HashDataUtils {
	var res models.HashDataUtils

	file, err := os.Open(path)
	if err != nil {
		errors.CheckErr(err)
		return models.HashDataUtils{}
	}

	defer file.Close()

	result, err := s.hasher.Hash(file)

	if err != nil {
		errors.CheckErr(err)
		return models.HashDataUtils{}
	}

	res.FileName = filepath.Base(path)
	res.Checksum = result
	res.FilePath = file.Name()
	res.Algorithm = s.alg

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
func (s *HashService) CallFunction(helpPath bool, dirPath string, getData bool, getChangedData string, updDeleted string, podData models.PodInfo) {
	switch {
	case helpPath:
		flag.Usage = func() {
			fmt.Fprintf(os.Stderr, "Help with commands %s:\nUse one of the following commands:\n", os.Args[0])
			flag.VisitAll(func(f *flag.Flag) {
				fmt.Fprintf(os.Stderr, " flag:	-%v \n 	%v\n", f.Name, f.Usage)
			})
		}
		flag.Usage()
	case len(dirPath) > 0:
		s.SavingData(s.CheckSum(dirPath), podData)
	case getData:
		s.GetData()
	case len(getChangedData) > 0:
		s.GetChangedData(getChangedData)
	case len(updDeleted) > 0:
		s.UpdateDeletedStatus(updDeleted)
	default:
		log.Println("Error with flag, use '-h' flag for help ")
	}
}

//Getting data
func (s *HashService) GetData() ([]models.HashData, error) {
	data, err := s.repo.GetDataFromDB()

	if data == nil {
		log.Println("no data for output")
		return nil, err
	}

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for _, h := range data {
		if h.Algorithm == s.alg {
			fmt.Printf("File name: %s, Checksum: %s, Algorithm: %s\n", h.FileName, h.CheckSum, h.Algorithm)
		}
	}
	return data, nil
}

//Inserting data
func (s *HashService) PutData(res []models.HashDataUtils, podData models.PodInfo) error {
	var data []models.HashData
	var podStruct models.PodData
	podStruct.PodName = podData.PodName
	podStruct.ContainerName = podData.ContainerName
	podStruct.ImageName = podData.ImageName
	podStruct.CreationTime = podData.CreationTime
	for _, h := range res {
		var dat models.HashData
		dat.FileName = h.FileName
		dat.FilePath = h.FilePath
		dat.Algorithm = h.Algorithm
		dat.CheckSum = h.Checksum
		data = append(data, dat)
	}
	return s.repo.PutDataInDB(data, podStruct)
}

func (s *HashService) GetChangedData(dir string) (int, error) {

	var code int

	code = 0

	data, err := s.repo.GetDataByPathFromDB(s.alg)

	if data == nil {
		log.Fatalln("no data for output")
	}

	if err != nil {
		log.Fatalln(err)
	}

	secondData := s.CheckSum(dir)
	for _, h := range data {
		for _, c := range secondData {
			switch {
			case h.FileName == c.FileName:
				if h.FilePath != c.FilePath || h.CheckSum != c.Checksum {
					fmt.Printf("something with this file: %s\n", c.FileName)
					code = 1
				}
			case h.FilePath == c.FilePath:
				if h.FileName != c.FileName || h.CheckSum != c.Checksum {
					fmt.Printf("something with this file by this path: %s\n", c.FilePath)
					code = 1
				}
			case h.CheckSum == c.Checksum:
				if h.FileName != c.FileName || h.FilePath != c.FilePath {
					fmt.Printf("something with this file by this checksum: %s\n", c.Checksum)
					code = 1
				}
			}
		}
	}
	if code == 0 {
		fmt.Println("files not changed, check successful")
	}

	return code, nil
}

func (s *HashService) UpdateDeletedStatus(dir string) error {
	var results []models.ChangedHashes
	var result models.ChangedHashes
	databaseData, err := s.repo.GetDataByPathFromDB(s.alg)

	if databaseData == nil {
		log.Println("no data for output")
		return err
	}

	if err != nil {
		fmt.Println(err)
	}

	secondData := s.CheckSum(dir)

	sm := make(map[string]struct{}, len(secondData))
	for _, n := range secondData {
		sm[n.FilePath] = struct{}{}
	}

	for _, n := range databaseData {
		if _, ok := sm[n.FilePath]; !ok {
			result.FilePath = n.FilePath
			result.Algorithm = n.Algorithm
			results = append(results, result)
		}
	}

	var data []models.HashData
	var dat models.HashData
	for _, h := range results {
		dat.FilePath = h.FilePath
		dat.Algorithm = h.Algorithm
		data = append(data, dat)
	}
	if data != nil {
		for _, h := range data {
			fmt.Printf("This file was deleted: file: %s, algorithm: %s\n", h.FilePath, h.Algorithm)
		}
	} else {
		fmt.Println("No deleted files founded")
	}
	s.repo.UpdateDeletedStatusInDB(data)
	return nil
}
