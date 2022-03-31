package errors

import (
	"log"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatalln("Error occurred: ", err)
	}
}
