package main

import (
	"fmt"
	v1beta1 "k8s.io/api/events/v1beta1"
	//	"context"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var namespace = "demo-test"
var podname = "hello-world-v2-99"

func main() {
	swatch()
}

func swatch()  {
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

	timeout := int64(90000)
	listOptions := meta_v1.ListOptions{
		TimeoutSeconds: &timeout,
		//LabelSelector:fmt.Sprintf("app=%s",podname),
	}

	byteReader, err :=  clientset.Events().Events(namespace).Watch(listOptions)

	//byteReader, err := clientset.CoreV1().Pods(namespace).Watch(listOptions)

	for {
		select {
		case event, ok := <-byteReader.ResultChan():
			if !ok {
				fmt.Println("WatchPod resultChan: ", ok)
				return
			}
			fmt.Println(event.Type)
			e := event.Object.(*v1beta1.Event)
			fmt.Println("event",e)
		}
	}

}