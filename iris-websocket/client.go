package main

import (
	"fmt"
	xwebsocket "golang.org/x/net/websocket"
)

func main() {
	//ws 服务端地址
	///websocket/status
	//var url = "ws://localhost:8080/websocket"
	//var url = "ws://localhost:8081/webSocket/appLogs?namespace=demo&&name=hello-world&&new=true"
	//var url =  "ws://localhost:8080/socket"
	var url = "ws://yxb.hl.10086.cn/hidevopsio/hiadmin/webSocket/buildLogs?namespace=demo&&name=hello-world"
	//var url = "ws://nginx-gateway.app.vpclub.io/demo/websocket/webSocket/buildLogs?namespace=moses&&name=merchants-provider"
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

//	fmt.Println(WS.HeaderReader())
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


//Head : map[X-Real-Ip:[10.128.0.1] X-Forwarded-Host:[hiadmin-demo.app.vpclub.io] X-Forwarded-For:[116.24.64.145, 192.168.1.50, 10.128.0.1 172.16.13.33] Origin:[http://localhost/ ] Sec-Websocket-Version:[13] X-Forwarded-Port:[80] X-Forwarded-Proto:[http] Forwarded:[for=172.16.13.33;host=hiadmin-demo.app.vpclub.io;proto=http] Sec-Websocket-Key:[iUkXQl2F417pFvKVTpbEFA==]]
//Head : map[X-Forwarded-Host:[hiadmin-demo.app.vpclub.io] Upgrade:[websocket] X-Forwarded-For:[116.24.64.145, 192.168.1.50, 10.128.0.1 172.16.13.33] X-Real-Ip:[10.128.0.1] Sec-Websocket-Version:[13] Connection:[Upgrade] X-Forwarded-Port:[80] X-Forwarded-Proto:[http] Forwarded:[for=172.16.13.33;host=hiadmin-demo.app.vpclub.io;proto=http] Sec-Websocket-Key:[iUkXQl2F417pFvKVTpbEFA==] Origin:[http://localhost/ ]]


//Head : map[X-Forwarded-Port:[80]
// X-Forwarded-Proto:[http]
// Upgrade:[websocket]
// X-Forwarded-For:[172.16.1.51, 10.128.0.1 172.16.13.144]
// Sec-Websocket-Key:[tMv9MGw5PsPIByjpD2+StA==]
// Sec-Websocket-Version:[13]
// X-Forwarded-Host:[hiadmin-demo.app.vpclub.io]
// Connection:[upgrade]
// X-Real-Ip:[10.128.0.1]
// Forwarded:[for=172.16.13.144;host=hiadmin-demo.app.vpclub.io;proto=http]]


//Head : map[Sec-Websocket-Version:[13]
// X-Forwarded-Host:[hiadmin-demo.app.vpclub.io]
// X-Forwarded-Proto:[http]
// X-Forwarded-For:[116.24.64.145, 192.168.1.50, 10.128.0.1 172.16.13.33]
// X-Real-Ip:[10.128.0.1]
// Sec-Websocket-Key:[712F5qDMD/JVmG/lMSF+Mg==]
// Origin:[http://localserver/ ]
// X-Forwarded-Port:[80]
// Forwarded:[for=172.16.13.33;host=hiadmin-demo.app.vpclub.io;proto=http]]



//Head : map[Upgrade:[websocket websocket]
// X-Forwarded-For:[172.16.1.51, 10.128.0.1 172.16.13.144]
// Sec-Websocket-Key:[tMv9MGw5PsPIByjpD2+StA==]
// Sec-Websocket-Version:[13]
// X-Forwarded-Host:[hiadmin-demo.app.vpclub.io]
// X-Forwarded-Port:[80]
// X-Forwarded-Proto:[http]
// Connection:[upgrade Upgrade]
// X-Real-Ip:[10.128.0.1]
// Forwarded:[for=172.16.13.144;host=hiadmin-demo.app.vpclub.io;proto=http]]

//========
//Head : map[Sec-Websocket-Version:[13] X-Forwarded-Host:[hiadmin-demo.app.vpclub.io] X-Forwarded-Proto:[http] X-Forwarded-For:[116.24.64.145, 192.168.1.50, 10.128.0.1 172.16.13.33] X-Real-Ip:[10.128.0.1] Sec-Websocket-Key:[712F5qDMD/JVmG/lMSF+Mg==] Origin:[http://localserver/ ] X-Forwarded-Port:[80] Forwarded:[for=172.16.13.33;host=hiadmin-demo.app.vpclub.io;proto=http]]
//Head : map[Sec-Websocket-Key:[712F5qDMD/JVmG/lMSF+Mg==] Sec-Websocket-Version:[13] Upgrade:[websocket] X-Forwarded-Host:[hiadmin-demo.app.vpclub.io] X-Forwarded-Proto:[http] X-Forwarded-For:[116.24.64.145, 192.168.1.50, 10.128.0.1 172.16.13.33] X-Real-Ip:[10.128.0.1] Forwarded:[for=172.16.13.33;host=hiadmin-demo.app.vpclub.io;proto=http] Connection:[Upgrade] Origin:[http://localserver/ ] X-Forwarded-Port:[80]]
