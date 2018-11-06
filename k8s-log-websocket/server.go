// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var podname string = "hello-world-36-d98fc675d-p85xz"
var namespace string = "demo"

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options


func GetLog(namespac,name string) (*bufio.Reader,error) {
	kubeconfig := flag.String("kubeconfig", "/Users/wang/.kube/config", "(optional) absolute path to the kubeconfig file")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	ctx := context.TODO()
	til := int64(10)
	byteReader, err := clientset.CoreV1().Pods(namespac).
		GetLogs(name, &v1.PodLogOptions{Follow: true,TailLines:&til}).Context(ctx).Stream()
	if err != nil {
		fmt.Println("Error ",err)
		return nil ,err
	}
	return bufio.NewReader(byteReader),nil
}

var Reader *bufio.Reader

func init() {
	var err error
	Reader,err = GetLog(namespace,podname)
	if err != nil {
		log.Print("Error ",err)
		return
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		fmt.Println("message",message)

		//err = c.WriteMessage(mt, []byte(fmt.Sprintf("Logs in Namespace:%s PodName: %s","hidevopsio-alpha","hidevopsio-log-6b94b49dbc-xr27t")))
		err = nil
		for err == nil {
			str, err := Reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error ",err)
				break
			}
			fmt.Println(str)
			err = c.WriteMessage(mt, []byte(str))
			if err != nil {
				log.Println("write:", err)
				break
			}

			if err != nil {
				fmt.Println("Error ",err)
				break
			}
		}
		if err == io.EOF {
			fmt.Println("Error ",err)
			return
		}

	}
}

func home(w http.ResponseWriter, r *http.Request) {

	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")

}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}


var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <script>
        window.addEventListener("load", function(evt) {
            var oDiv = document.getElementById('float_banner');
            oDiv.style.position = 'fixed';
            oDiv.style.top = '20px';
            oDiv.style.left = '20px';
            oDiv.style.backgroundColor = "#0000FF";
            oDiv.style.color="#F8F8FF";
            var float = function(message) {
                var d = document.createElement("div");
                d.innerHTML = message;
                oDiv.appendChild(d);
            };


            var output = document.getElementById("output");
            var logtext = document.getElementById("logtext");
            logtext.style.backgroundColor = "#000000";
            logtext.style.color="#F8F8FF";
            //var ws;

            var print = function(message) {
                var d = document.createElement("div");
                d.innerHTML = message;
                output.appendChild(d);
            };

//            var ws = new WebSocket("ws://localhost:8080/echo");
            ws = new WebSocket("{{.}}");
			ws.send("visit");
            ws.onopen = function(evt) {
                console.log("Connection open ...");
                ws.send("Hello WebSockets!");
            };
            var i = 0;
            ws.onmessage = function(evt) {
                if (i=0) {
                    float(evt.data);
                    i=1
                }


                print(evt.data)
                //ws.close();
            };

            ws.onclose = function(evt) {
                console.log("Connection closed.");
            };

        });
    </script>
</head>
<div id="logtext">
<div id="float_banner"><p></p></div>
<div id="output"></div>
</div>
</html>
`))
