package main

import (
	"fmt"
	"github.com/kataras/iris/core/errors"
	"time"
)

func main() {

	s := make(chan string)
	go func() {
		//if err := SendStr(s);err != nil {
		//	close(s)
		//}

		for i:=0;i <= 10;i++ {
			time.Sleep(time.Second)
			str := string(i)
			s <- str
		}

	}()

	for {
		select {
		case str :=  <- s:
			fmt.Println(string(str))
			break
		case <-time.After(2 * time.Second):
			fmt.Println("END")
			return
		}
	}
}

func SendStr(str chan *string) error {

	for i:=0;i <= 10;i++ {
		time.Sleep(time.Second)
		s := string(i)
		str <- &s
	}
	return errors.New("err")

}