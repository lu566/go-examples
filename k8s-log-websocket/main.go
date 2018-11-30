package main

import (
	"flag"
	"fmt"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
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
//"app=hello-world,version=v2"

	w,err := clientset.CoreV1().Pods("demo-dev").Watch(meta_v1.ListOptions{LabelSelector:"app=hello-world,version=v2"})
	if err != nil {
		return
	}
	fmt.Println(w)
	//FieldSelector:"involvedObject.name=admin-v1-347-55chj"
	//FieldSelector:"involvedObject.name%3Dhello-world-v2-794c59bd46-k58pm%2CinvolvedObject.namespace%3Ddemo-dev"


	name := "hello-world-v7-74dbc96d58-tgzzd"//generateName: hello-world-v7-74dbc96d58-
	eventSelector, _ := fields.ParseSelector("involvedObject.name=" + name)
	//eventSelector, _ := fields.ParseSelector("involvedObject.kind=Pod")
	e,err := clientset.CoreV1().Events("demo-test").Watch(meta_v1.ListOptions{
		FieldSelector:eventSelector.String(),
		//LabelSelector:fmt.Sprintf("app=%s","hello-world"),

		Watch:true,
	})
	if err != nil {
		fmt.Println("Error",err)
		return
	}


	fmt.Println(eventSelector.String())

	eventList,err := clientset.CoreV1().Events("demo-test").List(meta_v1.ListOptions{
		//FieldSelector:eventSelector.String(),
		LabelSelector:fmt.Sprintf("app=%s,version=%s","hello-world","v7"),
		//Watch:true,
	})
	fmt.Println(len(eventList.Items),eventList)

	for {
		select {
		//case event := <-w.ResultChan():
		//
		//	pod := event.Object.(*corev1.Pod)
		//	fmt.Println(event.Type,"===",pod.Status.Phase)

		case event := <-e.ResultChan():
//			fmt.Println("Event",event.Object.(*v1.Event))
			e := event.Object.(*v1.Event)
			fmt.Println(e.Message)
		}
	}
}
