package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sha256-sum/services"
)

var (
	filePath string
	helpPath bool
	dirPath  string
)

func init() {
	flag.StringVar(&filePath, "p", "", "Hash through file path")
	flag.StringVar(&dirPath, "d", "", "Hash of all files through directory path")
	flag.BoolVar(&helpPath, "h", false, "info")
	flag.Parse()
}

func main() {
	switch {
	case helpPath:
		flag.Usage = func() {
			fmt.Fprintf(os.Stderr, "Help with commands %s:\nUse one of the following commands:\n", os.Args[0])
			flag.VisitAll(func(f *flag.Flag) {
				fmt.Fprintf(os.Stderr, " flag 	-%v \n  	%v\n", f.Name, f.Usage)
			})
		}
		flag.Usage()
	case len(filePath) > 0:
		fmt.Println(services.HashOfFile(filePath))
	case len(dirPath) > 0:
		services.CheckSum(dirPath)
	default:
		log.Println("Error with flag, use '-h' flag for help ")
	}
}
