package main

import (
	"fmt"
	"os"
	"net"
	"io"
)

func SendFile(path string, connect net.Conn){

	file, oerr :=os.Open(path)
	if oerr !=nil{
		fmt.Println("Open", oerr)
		return
	}
	defer file.Close()
	buff := make([]byte,1024*4)
	for{
		size, rerr := file.Read(buff)
		if rerr != nil{
			if rerr == io.EOF{
				fmt.Println("EOF",rerr)

			}else{
				fmt.Println("Read:", rerr)
			}
			return
		}
		connect.Write(buff[:size])
	}
}

func main(){

	fmt.Print("请输入需要传输的文件路径:")
	var path string
	fmt.Scan(&path)

	info, serr :=os.Stat(path)
	if serr !=nil{
		fmt.Println("Stat", serr)
		return
	}

	connect, derr :=net.Dial("tcp","127.0.0.1:10240")
	if derr !=nil{
		fmt.Println("Dial", derr)
		return
	}

	_,werr:= connect.Write([]byte(info.Name()))
	if werr!=nil{
		fmt.Println("Write",werr)
		return
	}

	buff := make([]byte,4096)
	size,rerr := connect.Read(buff)
	if rerr!=nil{
		fmt.Println("Read",rerr)
		return
	}

	if "ok" == string(buff[:size]){
		SendFile(path, connect)
	}

	defer connect.Close()
}