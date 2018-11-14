package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	go func() {
		for i:=0;i<=100;i++ {
			fmt.Println(i)
			time.Sleep(time.Second)
		}
	}()

	for i:=0;i<=100;i++ {

		fmt.Println(os.Stdout.Name())

		time.Sleep(time.Second)
	}



}