package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

type WebsocketHandler struct {
	m *melody.Melody
}

func (wh *WebsocketHandler) NewWebsocketHandler(mel *melody.Melody) *WebsocketHandler {

	return &WebsocketHandler{
		m: mel,
	}
}

func (wh *WebsocketHandler) RegisterRoutes(r gin.IRouter) {

	wsApi := r.Group("/ws")

	wsApi.GET("/")
}
