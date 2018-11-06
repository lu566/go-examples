package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
	"time"
)

func main() {
	ServerLoop()
}

func ServerLoop() {
	app := iris.New()
	ws := websocket.New(websocket.Config{})
	app.Get("/socket", ws.Handler())
	ws.OnConnection(OnConnect)
	app.Run(iris.Addr(":8080"))
}

// OnConnect handles incoming websocket connection
func OnConnect(c websocket.Connection) {
	//Join 线程隔离


	//Join将此连接注册到房间，如果它不存在则会创建一个新房间。一个房间可以有一个或多个连接。一个连接可以连接到许多房间。所有连接都自动连接到由“ID”指定的房间。
	c.Join("room1")
	// Leave从房间中删除此连接条目//如果连接实际上已离开特定房间，则返回true。
	defer c.Leave("room1")

	//获取路径中的query数据
	fmt.Println("namespace:",c.Context().FormValue("namespace"))
	fmt.Println("nam:",c.Context().FormValue("name"))

	//收到消息的事件。bytes=收到客户端发送的消息
	c.OnMessage(func(bytes []byte) {
		fmt.Println("client:",string(bytes))

		//循环向连接的客户端发送数据
		var i int
		for {
			i++
			c.EmitMessage([]byte(fmt.Sprintf("=============%d==========",i)))   // ok works too
			time.Sleep(time.Second)
		}
	})
}