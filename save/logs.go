package aggregate

import (
	"bufio"
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hioak/starter/kube"
	"github.com/hidevopsio/mioclient/starter/mio"
	"github.com/kataras/iris/core/errors"
	"k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

//go:generate mockgen -destination mock_logs.go -package aggregate hidevops.io/hiadmin/pkg/console/aggregate KubeClient

type KubeClient interface {
	GetBuildConfigLastVersion(namespace, name string) (int,error)
	GetPodNameBylabel(namespace, name string,LastVersion int)(string,error)
	WatchPodStatus(namespace,name string,podMessage chan PodMessage) error
	GetLogs(namespace, name string,tail int64) (*bufio.Reader,error)
}

type KubeClientImpl struct {
	KubeClient
	pod 		 *kube.Pod
	buildConfig  *mio.BuildConfig
}

func init() {
	app.Register(newLogOutConfig)
}

func newLogOutConfig(pod *kube.Pod ,buildConfig *mio.BuildConfig) KubeClient {
	return &KubeClientImpl{
		pod:pod,
		buildConfig:buildConfig,
	}
}

type PodMessage struct {
	Message string
	IsEnd   bool
}

//获取 buildconfig 的lastVersion 字段信息
func (l *KubeClientImpl)GetBuildConfigLastVersion(namespace, name string) (int,error) {
	bc,err := l.buildConfig.Get(name,namespace)
	if err  != nil {
		fmt.Println("Error",err)
		return 0,err
	}
	return bc.Status.LastVersion,nil
}

//根据标签信息获取pod name
func (l *KubeClientImpl)GetPodNameBylabel(namespace, name string,LastVersion int)(string,error){

	labelSelector := fmt.Sprintf("app=%s-%d",name,LastVersion)
	podList,err := l.pod.GetPodList(namespace,meta_v1.ListOptions{
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

//获取pod Stream 信息，并且返回
func (l *KubeClientImpl)GetLogs(namespace, name string,tail int64) (*bufio.Reader,error) {
	podLogOptions := &v1.PodLogOptions{Follow: true}
	if tail != 0 {
		podLogOptions.TailLines = &tail
	}

	request, err := l.pod.GetPodLogs(namespace,name,podLogOptions)
	if err != nil {
		fmt.Println("Error ",err)
		return nil,err
	}

	byteReader,err := request.Stream()
	if err != nil {
		fmt.Println("Error ",err)
		return nil,err
	}

	reader := bufio.NewReader(byteReader)
	return reader,nil
}

//监听pod的状态，等待pod 正常或者失败 ，超时时间为10分钟
func (l *KubeClientImpl)WatchPodStatus(namespace,name string,podMessage chan PodMessage) error {
	timeout := int64(3600)

	listOptions := meta_v1.ListOptions{
		TimeoutSeconds: &timeout,
	}

	w, err := l.pod.Watch(listOptions,namespace,name)
	if err != nil {
		return err
	}

	startTime := time.Now().Second()
	for {
		select {
		case event, ok := <-w.ResultChan():
			if !ok {
				fmt.Println("WatchPod resultChan: ", ok)
				return nil
			}
			pod := event.Object.(*corev1.Pod)
			if pod.Status.Phase == corev1.PodPending {
				podMessage <- PodMessage{Message: fmt.Sprintf("[EVENT] type %s, Pod %s,status %s", pod.Name, event.Type, pod.Status.Phase), IsEnd: false}

				if time.Now().Second() - startTime >= 600 {
					podMessage <- PodMessage{Message: fmt.Sprintf("[EVENT] type %s, Pod %s,status %s", pod.Name, event.Type, pod.Status.Phase), IsEnd: true}
					return errors.New("Pod query timeout 10 minutes")
				}
				continue
			}

			if pod.Status.Phase == corev1.PodRunning {
				podMessage <- PodMessage{Message: fmt.Sprintf("[EVENT] type %s, Pod %s,status %s", pod.Name, event.Type, pod.Status.Phase), IsEnd: true}
				fmt.Printf("pod %s has been running", pod.Name)
				return nil
			}else {
				podMessage <- PodMessage{Message: fmt.Sprintf("[EVENT] type %s, Pod %s,status %s", pod.Name, event.Type, pod.Status.Phase), IsEnd: true}

				return fmt.Errorf("Pod type %s status %s",string(event.Type),pod.Status.Phase)
			}
		}
	}

	return nil
}