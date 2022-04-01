package main

import (
	"flag"
	"fmt"
	"log"
	"sha256-sum/services"
)

var filePath string

func init() {
	flag.StringVar(&filePath, "p", "", "file path")
	flag.Parse()
}

func main() {
	if len(filePath) > 0 {
		fmt.Println(services.HashOfFile(filePath))
	} else {
		log.Println("Error with flag, use '-p' flag ")
	}
}
