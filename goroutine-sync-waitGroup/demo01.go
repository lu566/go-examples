package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	//go func() {
	//	wg.Wait()
	//}()

	go func() {
		time.Sleep(time.Second * 5)
		wg.Done()
	}()

	go func() {


		i := 0
		for  {

			i++
			fmt.Println(i)
			time.Sleep(time.Second)
		}

	}()


	time.Sleep(time.Second * 1000000)




}