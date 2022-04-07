package main

import (
	"flag"
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
	services.CatchStopSignal()
	services.CallFunction(filePath, helpPath, dirPath, hashAlg)
}
