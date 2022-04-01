package main

import (
	"flag"
	"fmt"
	"log"
	"sha256-sum/services"
)

var filePath string
var dirPath string

func init() {
	flag.StringVar(&filePath, "p", "", "file path")
	flag.StringVar(&dirPath, "d", "", "dir path")
	flag.Parse()
}

func main() {
	switch {
	case len(filePath) > 0:
		fmt.Println(services.HashOfFile(filePath))
	case len(dirPath) > 0:
		fmt.Println(services.HashOfDir(dirPath))
	default:
		log.Println("Error with flag, use '-p' flag ")
	}
}
