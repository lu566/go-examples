package wsservice

import (
	"errors"
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/websocket"
	"hidevops.io/hiadmin/pkg/console/aggregate"
	"io"
	"time"
)

type countHandler struct {
	connection websocket.Connection
	kubeClient aggregate.KubeClient
}

type CountHandlerConstructor interface{}

func NewCountHandlerConstructor() CountHandlerConstructor {
	return func(connection websocket.Connection) websocket.Handler {
		return &countHandler{connection: connection}
	}
}

func init() {
	app.Register(new(CountHandlerConstructor), NewCountHandlerConstructor)
}

func (h *countHandler) OnMessage(data []byte) {
	message := string(data)
	log.Debugf("client: %v", message)

	namespace := h.connection.Context().FormValue("namespace")
	name := h.connection.Context().FormValue("name")

	//获取 buildConfig lastVersion
	lastVersion,err := h.kubeClient.GetBuildConfigLastVersion(namespace,name)
	if err != nil {
		return
	}

	//根据lastVersion 拼写标签信息 获取 podName
	podName,err := h.kubeClient.GetPodNameBylabel(namespace,name,lastVersion)
	if err != nil {
		return
	}

	//使用协程实时向客户端发送状态
	podMessage := make(chan aggregate.PodMessage)
	go h.kubeClient.WatchPodStatus(namespace,podName, podMessage)
	for {
		select {
		case <- time.After(10 * time.Minute):
			fmt.Println(errors.New("Pod query timeout 10 minutes"))
			return

		case m := <- podMessage:
			//TODO WS send message
			if err := h.connection.EmitMessage([]byte(m.Message));err != nil {
				fmt.Println("Error ",err)
				break
			}
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

		fmt.Println(str)
		//TODO WS send message
		if err := h.connection.EmitMessage([]byte(string(str)));err != nil {
			fmt.Println("Error ",err)
			break
		}
	}

	if err == io.EOF {
		fmt.Println("Error ",err)
		return
	}

}

func (h *countHandler) OnDisconnect() {
	log.Debugf("Connection with ID: %v has been disconnected!", h.connection.ID())
}
