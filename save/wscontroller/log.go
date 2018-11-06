package wscontroller

import (
	"github.com/hidevopsio/hiboot/examples/web/websocket/service"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/starter/websocket"
	"hidevops.io/hiadmin/pkg/console/wsservice"
)

type websocketController struct {
	at.RestController
	connectionFunc           websocket.ConnectionFunc
	countHandlerConstructor  wsservice.CountHandlerConstructor

}

func newWebsocketController(connectionFunc websocket.ConnectionFunc,
	countHandlerConstructor service.CountHandlerConstructor) *websocketController {
	c := &websocketController{
		connectionFunc:           connectionFunc,
		countHandlerConstructor:  countHandlerConstructor,
	}
	return c
}

func init() {
	app.Register(newWebsocketController)
}

func (c *websocketController) GetLogs(ctx *web.Context) {
	c.connectionFunc(ctx, c.countHandlerConstructor.(func(websocket.Connection) websocket.Handler))
}
