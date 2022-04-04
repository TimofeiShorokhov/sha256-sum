package main

import (
	"flag"
	"fmt"
	"log"
	"sha256-sum/services"
	"sync"
)

var (
	filePath string
	helpPath string
	dirPath  string
	wg       sync.WaitGroup
)

func init() {
	flag.StringVar(&filePath, "p", "", "file path")
	flag.StringVar(&dirPath, "d", "", "dir path")
	flag.StringVar(&helpPath, "h", "", "")
	flag.Parse()
}

func main() {
	wg.Add(1)
	switch {
	case len(filePath) > 0:
		fmt.Println(services.HashOfFile(filePath, &wg))
	case len(dirPath) > 0:
		go services.HashOfDir(dirPath, &wg)
	case len(helpPath) > 0:
		services.PrintHelp(helpPath, &wg)
	default:
		log.Println("Error with flag, use '-h' flag with '--help' argument for help ")
	}
	wg.Wait()
}
