// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	mio_v1alpha1 "github.com/hidevopsio/mioclient/pkg/client/clientset/versioned/typed/mio/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	restclient "k8s.io/client-go/rest"
	"bufio"
	"context"
	"flag"
	"fmt"
	"html/template"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"net/http"
	"io"
	"github.com/gorilla/websocket"
)

var podname string = "hello-world-v2-95-7b69bbbc67-b546j"
var namespace string = "demo"

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

var config *restclient.Config
func init(){

	kubeconfig := flag.String("kubeconfig", "/Users/wang/.kube/config", "(optional) absolute path to the kubeconfig file")
	flag.Parse()
	var err error
	config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
}

func GetBuildConfigLastVersion(namespace, name string) (int,error) {

	mioClient,err :=mio_v1alpha1.NewForConfig(config)
	if err  != nil {
		fmt.Println("Error",err)
		return 0,err
	}

	bc,err := mioClient.BuildConfigs(namespace).Get(name,meta_v1.GetOptions{})
	if err  != nil {
		fmt.Println("Error",err)
		return 0,err
	}
	return bc.Status.LastVersion,nil
}

func GetPodNameBylabel(namespace, name string,LastVersion int)(string,error){
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "",err
	}

	labelSelector := fmt.Sprintf("app=%s-%d",name,LastVersion)
	podList,err := clientset.CoreV1().Pods(namespace).List(meta_v1.ListOptions{
		LabelSelector:labelSelector,
	})
	if err != nil {
		return "",err
	}
	for _,pod := range podList.Items {
		fmt.Println("INFO found pod name :",pod.Name)
	}

	if len(podList.Items)>1 {
		return "",fmt.Errorf("The label %s matching pod should have only one",labelSelector)
	}else if len(podList.Items) == 0 {
		return "",fmt.Errorf("The label %s find pod failed",labelSelector)
	}
	return podList.Items[0].Name,nil
}

func GetLogs(namespace, name string,tail int) (*bufio.Reader,error) {

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	ctx := context.TODO()
//	lin := int64(10)
	byteReader, err := clientset.CoreV1().Pods(namespace).
		GetLogs(name, &v1.PodLogOptions{Follow: true,}).Context(ctx).Stream()
	if err != nil {
		fmt.Println("Error ",err)
		return nil,err
	}

	reader := bufio.NewReader(byteReader)
	return reader,nil

	//err = nil
	//for err == nil {
	//	str, err := reader.ReadString('\n')
	//	fmt.Println("--",str)
	//	go func() {
	//		types.Message <- str
	//	}()
	//	if err != nil {
	//		fmt.Println("Error ",err)
	//		break
	//	}
	//	if err != nil {
	//		fmt.Println("Error ",err)
	//		break
	//	}
	//}
	//if err == io.EOF {
	//	fmt.Println("Error ",err)
	//	return err
	//}
	//return nil
}

func WatchPodStatus(namespace,labelSelector string) error {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	timeout := int64(9000)

	listOptions := meta_v1.ListOptions{
		TimeoutSeconds: &timeout,
		LabelSelector:labelSelector,
	}

	w, err := clientset.CoreV1().Pods(namespace).Watch(listOptions)
	if err != nil {
		return err
	}

		for {
			select {
			case event, ok := <-w.ResultChan():
				if !ok {
					fmt.Println("WatchPod resultChan: ", ok)
					return nil
				}
				switch event.Type {
				case watch.Added:
					pod := event.Object.(*corev1.Pod)
					if pod.Status.Phase == "Running" {
						fmt.Println("pod %s has been running", pod.Name)
						return nil
					}
				case watch.Modified:
					pod := event.Object.(*corev1.Pod)
					if pod.Status.Phase == "Running" {
						fmt.Println("pod %s has been running", pod.Name)
						return nil
					}
				case watch.Deleted:
					fmt.Println("Deleted: ", event.Object)
					return fmt.Errorf("pod has been removed")
				}
			}
		}

	return nil
}


func logs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("logs")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	namespace := r.FormValue("namespace")
	name := r.FormValue("name")

	if name == "" || namespace == "" {
		log.Println(fmt.Errorf("namespace or name cannot be empty"))
		return
	}

	lastVersion,err := GetBuildConfigLastVersion(namespace,name)
	if err != nil {
		log.Println("Error ",err)
		return
	}
	fmt.Println("lastversion",lastVersion)

	podName,err := GetPodNameBylabel(namespace,name,lastVersion)
	if err != nil {
		log.Println("Error ",err)
		return
	}

	labelSelector := fmt.Sprintf("app=%s",name)
	if err := WatchPodStatus(namespace,labelSelector);err != nil {
		log.Println("Error ",err)
		return
	}


	reader,err := GetLogs(namespace,podName,0)
	if err != nil {
		log.Println("Error ",err)
		return
	}

	//mt, message, err := c.ReadMessage()
	//if err != nil {
	//	log.Println("read:", err)
	//	return
	//}
	//fmt.Println(message)
	//
	//if err := c.WriteMessage(mt,[]byte("zzzzzzzzzzzz"));err != nil {
	//	log.Println("Error ",err)
	//	return
	//}

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		fmt.Printf("message %s",message)

		//err = c.WriteMessage(mt, []byte(fmt.Sprintf("Logs in Namespace:%s PodName: %s","hidevopsio-alpha","hidevopsio-log-6b94b49dbc-xr27t")))
		err = nil
		for err == nil {
			str, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error ",err)
				break
			}
			fmt.Println("str",str)
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

	homeTemplate.Execute(w, "ws://"+r.Host+"/logs?namespace=demo&name=hello-world")

}

func main() {
	//test
	//
	//if err := WatchPodStatus("demo","hello-world-48",);err != nil {
	//	fmt.Println("Error ",err)
	//}
	//return


	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/logs", logs)
	http.HandleFunc("/", home)
	http.ListenAndServe(*addr, nil)
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
