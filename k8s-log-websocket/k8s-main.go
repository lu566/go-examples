package main

import (
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sync"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/api/core/v1"
	//"k8s.io/klog"
	corev1 "k8s.io/api/core/v1"
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

	app := "hello-world"
	version := "v7"
	namespace := "demo-test"

	podWatch,err := clientset.CoreV1().Pods(namespace).Watch(metav1.ListOptions{
		LabelSelector:fmt.Sprintf("app=%s,version=%s",app,version),
	})

	if err != nil {
		fmt.Println("Error ",err)
		return
	}

	for  {
		select {

		case p := <-podWatch.ResultChan():
			switch p.Type {
			case watch.Added:
				pod := p.Object.(*corev1.Pod)
				fmt.Println("Pod",pod.Name)


				ch,err := WatchPodstatus(pod.Name, pod.Namespace, make(chan struct{}),clientset)
				if err != nil {
					fmt.Println("Error", err)
				}

				waitForPod(pod,ch)

			}

		}

	}

}


// WatchPod watches a pod and events related to it. It sends pod updates and events over the returned channel
// It will continue until stopChannel is closed
func WatchPodstatus(name, namespace string, stopChannel chan struct{},c *kubernetes.Clientset) (<-chan watch.Event, error) {
	podSelector, err := fields.ParseSelector("metadata.name=" + name)
	if err != nil {
		return nil, err
	}
	options := metav1.ListOptions{
		FieldSelector: podSelector.String(),
		Watch:         true,
	}
	podWatch, err := c.CoreV1().Pods(namespace).Watch(options)
	if err != nil {
		return nil, err
	}

	eventSelector, _ := fields.ParseSelector("involvedObject.name=" + name)
	eventWatch, err := c.CoreV1().Events(namespace).Watch(metav1.ListOptions{
		FieldSelector: eventSelector.String(),
		Watch:         true,
	})
	if err != nil {
		podWatch.Stop()
		return nil, err
	}

	eventWatchList, err := c.CoreV1().Events(namespace).List(metav1.ListOptions{
		FieldSelector: eventSelector.String(),
	})
	if err != nil {
		return nil,err
	}

	fmt.Println(len(eventWatchList.Items),eventWatchList)

	eventCh := make(chan watch.Event, 30)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer close(eventCh)
		wg.Wait()
	}()

	go func() {
		defer eventWatch.Stop()
		defer wg.Done()
		for {
			select {
			case _ = <-stopChannel:
				return
			case eventEvent, ok := <-eventWatch.ResultChan():
				if !ok {
					return
				} else {
					eventCh <- eventEvent
					fmt.Println(eventEvent.Type,eventEvent.Object)
				}
			}
		}
	}()

	go func() {
		defer podWatch.Stop()
		defer wg.Done()
		for {
			select {
			case <-stopChannel:
				return

			case podEvent, ok := <-podWatch.ResultChan():
				if !ok {
					return
				} else {
					eventCh <- podEvent
					//fmt.Println(podEvent.Type,podEvent.Object)
				}
			}
		}
	}()

	return eventCh, nil
}



// waitForPod watches the pod it until it finishes and send all events on the
// pod to the PV.
func waitForPod(pod *v1.Pod, podCh <-chan watch.Event) error {
	for {
		event, ok := <-podCh
		if !ok {
			return fmt.Errorf("recycler pod %q watch channel had been closed", pod.Name)
		}
		switch event.Object.(type) {
		case *v1.Pod:
			// POD changed
			pod := event.Object.(*v1.Pod)
			//klog.V(4).Infof("recycler pod update received: %s %s/%s %s", event.Type, pod.Namespace, pod.Name, pod.Status.Phase)
			switch event.Type {
			case watch.Added, watch.Modified:
				if pod.Status.Phase == v1.PodSucceeded {
					// Recycle succeeded.
					return nil
				}
				if pod.Status.Phase == v1.PodFailed {
					if pod.Status.Message != "" {
						return fmt.Errorf(pod.Status.Message)
					} else {
						return fmt.Errorf("pod failed, pod.Status.Message unknown.")
					}
				}

			case watch.Deleted:
				return fmt.Errorf("recycler pod was deleted")

			case watch.Error:
				return fmt.Errorf("recycler pod watcher failed")
			}

		case *v1.Event:
			// Event received
			podEvent := event.Object.(*v1.Event)
			fmt.Println(podEvent.Message)			//klog.V(4).Infof("recycler event received: %s %s/%s %s/%s %s", event.Type, podEvent.Namespace, podEvent.Name, podEvent.InvolvedObject.Namespace, podEvent.InvolvedObject.Name, podEvent.Message)
			if event.Type == watch.Added {
				//recyclerClient.Event(podEvent.Type, podEvent.Message)
			}
		}
	}
}