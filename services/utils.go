package services

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sha256-sum/errors"
)

func HashOfFile(path string, hashAlg string) string {

	file, err := os.Open(path)
	if err != nil {
		errors.CheckErr(err)
		return ""
	}

	defer file.Close()

	var checkSum interface{}

	switch hashAlg {
	case "md5":
		hash := md5.New()
		_, err = io.Copy(hash, file)
		checkSum = hash.Sum(nil)
	case "512":
		hash := sha512.New()
		_, err = io.Copy(hash, file)
		checkSum = hash.Sum(nil)
	default:
		hash := sha256.New()
		_, err = io.Copy(hash, file)
		checkSum = hash.Sum(nil)
	}
	if err != nil {
		errors.CheckErr(err)
		return ""
	}
	res := fmt.Sprintf("%x %s", checkSum, filepath.Base(path))
	return res
}

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

func CallFunction(filePath string, helpPath bool, dirPath string, hashAlg string) {
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
		CheckSum(dirPath, hashAlg)
	default:
		log.Println("Error with flag, use '-h' flag for help ")
	}
}
