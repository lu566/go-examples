package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("old",time.Now().Unix())
	currTime := time.Now().Unix()
	time.Sleep(time.Second * 3)
	fmt.Println("new",time.Now().Unix())
	fmt.Println(currTime - time.Now().Unix())
}