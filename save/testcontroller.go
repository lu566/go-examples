package controller

import (
	"errors"
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/model"
	"hidevops.io/hiadmin/pkg/console/aggregate"
	"io"
	"time"
)

type LogOutController struct {
	web.Controller
	kubeClient aggregate.KubeClient
}

func init() {
	web.RestController(newLogsController)
}

func newLogsController(kubeClient aggregate.KubeClient) *LogOutController {
	return &LogOutController{
		kubeClient:kubeClient,
	}
}

func (h *LogOutController) GetLogs() (r model.Response,err error) {
	var name,namespace string
	if name == ""{
		name = "hello-world"
	}
	if namespace == "" {
		namespace = "demo"
	}

	fmt.Println("namespace",namespace,"name",name)

	//获取 buildConfig lastVersion
	lastVersion,err := h.kubeClient.GetBuildConfigLastVersion(namespace,name)
	if err != nil {
		return
	}

	pipelineConfigClientVersion,err := h.kubeClient.GetPipelineConfigClientVersion(namespace,name)
	if err != nil {
		return
	}

	label := fmt.Sprintf("%s-%s-%d",name,pipelineConfigClientVersion,lastVersion)

	//根据lastVersion 拼写标签信息 获取 podName
	podName,err := h.kubeClient.GetPodNameBylabel(namespace,fmt.Sprintf("app=%s",label))
	if err != nil {
		return
	}

	//使用协程实时向客户端发送状态
	podMessage := make(chan aggregate.PodMessage)
	go h.kubeClient.WatchPodStatus(namespace,label, podMessage)
	for {
		select {
		case <- time.After(10 * time.Minute):
			fmt.Println(errors.New("Pod query timeout 10 minutes"))
			return

		case m := <- podMessage:
				fmt.Println(m.Message)
			if m.IsEnd{
				goto exit
			}
		}
	}
exit:
	//获取 pod 日志输出流
	reader,err := h.kubeClient.GetLogs(namespace,podName,0)
	if err != nil {
		return
	}

	//持续向客户端发送日志信息
	for err == nil {

		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error ",err)
			break
		}

		fmt.Print(str)
		//TODO WS send message
		//if err := h.connection.EmitMessage([]byte(string(str)));err != nil {
		//	fmt.Println("Error ",err)
		//	break
		//}
	}

	if err == io.EOF {
		fmt.Println("Error ",err)
		return
	}


	base := new(model.BaseResponse)
	base.SetData(struct {
		Status string
	}{Status:"OK"})
	return base, nil
}
