package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)
var podname string = "hidevopsio-log-6b94b49dbc-xr27t"
var namespace string = "hidevopsio-alpha"
func main() {
	if err := GetLogs(namespace,podname);err != nil {
		fmt.Println("Error",err)
	}
}

func GetLogs(namespace, name string) error {
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


	mm := int64(10)
	ctx := context.TODO()
	byteReader, err := clientset.CoreV1().Pods("moses").
		GetLogs("invoice-admin-v2-2-8457fdcd75-hjdqh", &v1.PodLogOptions{SinceSeconds:&mm}).Context(ctx).Stream()
	if err != nil {
		fmt.Println("Error ",err)
		return err
	}
	reader := bufio.NewReader(byteReader)
	//str, err := reader.ReadByte()
	//
	//fmt.Println(string(str))
	var buf bytes.Buffer
	out := bufio.NewWriter(&buf)
	_, err = io.Copy(out, reader)
	fmt.Println(buf.String())


	return nil
	//reader := bufio.NewReader(byteReader)
	//err = nil
	//for err == nil {
	//	str, err := reader.ReadString('\n')
	//	fmt.Println("--",str)
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
	return nil
}