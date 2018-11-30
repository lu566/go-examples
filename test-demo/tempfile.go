package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	buf := "Hello, World"
	file, err := ioutil.TempFile("", "tmpfile")
	if err != nil {
		panic(err)
	}

	fmt.Println(file.Name())
	fmt.Println(file.Chdir())

	defer os.Remove(file.Name())

	if _, err := file.Write([]byte(buf)); err != nil {
		panic(err)
	}

	fmt.Println(file.Name())

}