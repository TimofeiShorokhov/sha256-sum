package main

import (
	"flag"
	"log"
	"sha256-sum/services"
)

func main() {
	var filePath string

	flag.StringVar(&filePath, "p", "", "file path")
	flag.Parse()

	if len(filePath) > 0 {
		services.HashOfFile(filePath)
	} else {
		log.Println("Error with flag, use '-p' flag ")
	}
}
