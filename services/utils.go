package services

import (
	"flag"
	"fmt"
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

type HashDataUtils struct {
	FileName  string
	Checksum  string
	FilePath  string
	Algorithm string
}

type ChangedHashes struct {
	FileName    string
	OldChecksum string
	NewChecksum string
	FilePath    string
	Algorithm   string
}

//Generating checksum of file
func (s *HashService) HashOfFile(path string) HashDataUtils {
	var res HashDataUtils

	file, err := os.Open(path)
	if err != nil {
		errors.CheckErr(err)
		return HashDataUtils{}
	}

	defer file.Close()

	result, err := s.hasher.Hash(file)

	if err != nil {
		errors.CheckErr(err)
		return HashDataUtils{}
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
func (s *HashService) CallFunction(helpPath bool, dirPath string, getData bool, getChangedData string, updDeleted string) {
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
		s.SavingData(s.CheckSum(dirPath))
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
func (s *HashService) GetData() ([]repository.HashData, error) {
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
func (s *HashService) PutData(res []HashDataUtils) error {
	var data []repository.HashData
	for _, h := range res {
		var dat repository.HashData
		dat.FileName = h.FileName
		dat.FilePath = h.FilePath
		dat.Algorithm = h.Algorithm
		dat.CheckSum = h.Checksum
		data = append(data, dat)
	}
	return s.repo.PutDataInDB(data)
}

func (s *HashService) GetChangedData(dir string) error {
	var results []ChangedHashes
	var result ChangedHashes
	data, err := s.repo.GetDataByPathFromDB(dir, s.alg)

	if data == nil {
		log.Println("no data for output")
		return err
	}

	if err != nil {
		fmt.Println(err)
	}

	secondData := s.CheckSum(dir)
	for _, h := range data {
		for _, c := range secondData {
			if h.FileName == c.FileName && h.Algorithm == c.Algorithm && h.CheckSum != c.Checksum {
				result.FileName = h.FileName
				result.FilePath = h.FilePath
				result.Algorithm = h.Algorithm
				result.OldChecksum = h.CheckSum
				result.NewChecksum = c.Checksum
				results = append(results, result)
				fmt.Printf("Checksum of this file: %s, by this path: %s, was changed. Checksum in database: %s, new checksum: %s, algorithm: %s\n",
					result.FileName, result.FilePath, result.OldChecksum, result.NewChecksum, result.Algorithm)
			}
		}
	}

	if results == nil {
		fmt.Println("No changes found")
	}
	return nil
}

func (s *HashService) UpdateDeletedStatus(dir string) error {
	var results []ChangedHashes
	var result ChangedHashes
	databaseData, err := s.repo.GetDataByPathFromDB(dir, s.alg)

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

	var data []repository.HashData
	var dat repository.HashData
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
