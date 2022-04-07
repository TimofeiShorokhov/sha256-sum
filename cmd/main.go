package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sha256-sum/services"
)

var (
	filePath string
	helpPath bool
	dirPath  string
	hashAlg  string
)

func init() {
	flag.StringVar(&filePath, "p", "", "Hash through file path")
	flag.StringVar(&dirPath, "d", "", "Hash of all files through directory path")
	flag.StringVar(&hashAlg, "a", "", "md5, sha512, default: sha256")
	flag.BoolVar(&helpPath, "h", false, "info")
	flag.Parse()
}

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("stopped by user %d", sig)
			os.Exit(1)
		}
	}()

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
		fmt.Println(services.HashOfFile(filePath, hashAlg))
	case len(dirPath) > 0:
		services.CheckSum(dirPath, hashAlg)
	default:
		log.Println("Error with flag, use '-h' flag for help ")
	}
}
