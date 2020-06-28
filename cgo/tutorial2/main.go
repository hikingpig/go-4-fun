package main

/*
  #include "hello.c"
*/
import "C"
import (
	"errors"
	"log"
)

func main() {
	//Call to void function without params
	err := Hello()
	if err != nil {
		log.Fatal(err)
	}
}

func Hello() error {
	_, err := C.Hello() //We ignore first result as it is a void function
	if err != nil {
		return errors.New("error calling Hello function: " + err.Error())
	}

	return nil
}
