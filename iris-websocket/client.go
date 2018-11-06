package main

import (
	"fmt"
	xwebsocket "golang.org/x/net/websocket"
)

func main() {
	//ws 服务端地址
	var url = "ws://localhost:8081/websocket/buildLogs?namespace=demo&name=hello-world&verbose=true"
	ClientLoop(url)
}

// ClientLoop connects to websocket server
func ClientLoop(url string) error {
	//连接 ws 服务端
	WS, err := xwebsocket.Dial(url, "", "http://localhost/")
	if err != nil {
		fmt.Println("failed to connect websocket", err.Error())
		return err
	}
	defer func() {
		if WS != nil {
			WS.Close()
		}
	}()

	//向服务端发送消息
	WS.Write([]byte("hello service"))

	//循环接收服务端发送过来的消息
	var msg = make([]byte, 2048)
	for {
		if n, err := WS.Read(msg); err != nil {
			fmt.Println(err.Error())
			return err
		} else {
			fmt.Println(string(msg[:n]))
		}

		//WS.Write([]byte("hello service"))
	}
}
