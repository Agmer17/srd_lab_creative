package ws

import (
	"fmt"

	"github.com/Agmer17/srd_lab_creative/internal/shared/middleware"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

type WebsocketHandler struct {
	mel *melody.Melody
}

func NewWebsocketHandler(m *melody.Melody) *WebsocketHandler {

	return &WebsocketHandler{
		mel: m,
	}
}

func (wh *WebsocketHandler) HandleHandshakeRequest(c *gin.Context) {

	userId, ok := middleware.GetUserID(c)
	if !ok {
		fmt.Println("no user id found where the fuck do i get this shit tho")
		c.JSON(403, "no userid found ")
		return
	}
	c.Request.Header.Set("X-User-ID", userId.String())
	wh.mel.HandleRequest(c.Writer, c.Request)

}

func (wh *WebsocketHandler) RegisterRoutes(r gin.IRouter) {

	wsEndpoint := r.Group("/ws")
	wsEndpoint.Use(middleware.AuthMiddleware())
	wsEndpoint.GET("/", wh.HandleHandshakeRequest)

}
