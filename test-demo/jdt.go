package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//
	//for i:=0; i<30; i++ {
	//	fmt.Println(r.Intn(5))
	//}

	//spinner(time.Microsecond)
	//rnd := rand.Intn(5)
	//fmt.Println(rnd)
	xh()
}


func spinner(delay time.Duration) {

		endStr := "========================> [%100]"
		str    := ">------------------------"
		i := 0
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		   if i >= 100 {
			   i = 100
			   fmt.Printf("\r%s", endStr)
			   str = ">------------------------"
			   i = 0
			   //return
		   }
		   str = strings.Replace(str,">-","=>",1)
		   fmt.Printf("\r%s [%s%d]", str,"%",i)

			time.Sleep(time.Second)
	   	}


}

func xh()  {
	endStr := "========================>"
	str    := ">------------------------"

		i := 1
		for {
			i= i+4
			if i >= 100 {
				fmt.Printf("\r%s/n", endStr)
				i = 4
				str = ">------------------------"
				//return
			}
			str = strings.Replace(str,">-","=>",1)
			fmt.Printf("\r%s", str)
			time.Sleep(time.Second)
		}


}