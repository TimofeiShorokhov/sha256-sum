package errors

import (
	"log"
)

func CheckErr(err error) {
	log.Println("Error occurred: ", err)
}
