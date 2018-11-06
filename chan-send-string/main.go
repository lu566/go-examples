package main

import (
	"fmt"
	"time"
)

func fibonacci(c, quit chan int) {
	x, y := 1, 1
	for {
		select {
		case c <- x:
			x, y = y, x + y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

var str = make(chan string)

func main() {

	go func() {
		for {
			select {
			case a := <-str:
				println("aaaa",a)
			}
		}
	}()


	time.Sleep(time.Second*2)
	fmt.Println("sleep 2 second")
	str <- ""

	time.Sleep(time.Second*1)
	fmt.Println("sleep 1 second")
	str <- "ss"

	time.Sleep(time.Second*1)
}


