package helper

import (
	"fmt"
	"log"
)

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func PrintIfError(err error, str string) {
	if err != nil {
		println(str)
	}
}

func LogIfError(err error, str string) {
	if err != nil {
		log.Println(str)
	}
}

func LogFatalIfError(err error, str string) {
	if err != nil {
		log.Fatalf(str, err)
	}
}

func Print(str string) {
	fmt.Print(str)
}
