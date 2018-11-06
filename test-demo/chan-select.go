package main

import (
	"fmt"
	"github.com/kataras/iris/core/errors"
	"time"
)

type status struct {
	message string
	isEnd bool
}

func receiveMsg(){
	c1 := make(chan status)

	var err error
	go func(){
		err = write(c1)

	}()


	//print two message
	for {
		select {
		case msg1:= <-c1:
			fmt.Println(msg1)
			if msg1.isEnd == true {
				fmt.Println("END")
				goto exit
			}

		}
	}
	exit:


}
func write(c1 chan status)error{
	for i:=0; i < 3; i++ {
		time.Sleep(time.Second * 1)
		c1 <- status{message:"6454a654asdsad",isEnd:false}
	}
	c1 <- status{message:"6454a654asdsad",isEnd:true}
	return errors.New("err")
}


func main() {
	receiveMsg()
}