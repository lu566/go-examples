package main

import (
	"fmt"
	"time"
)

func main() {

	var b = make(chan bool)

	go func() {

		b <- true
		time.Sleep(time.Second*3)
		b <- false
	}()

	for i:=0;i<=10;i++ {
		time.Sleep(time.Second*1)
		select {
		case b:= <-b:
			fmt.Println(b)
		default:
		}

	}


}