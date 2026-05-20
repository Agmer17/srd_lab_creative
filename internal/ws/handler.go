package ws

import (
	"github.com/Agmer17/srd_lab_creative/pkg"
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
	token := c.Query("token")
	if token == "" {
		c.JSON(403, "no userid found ")
		return
	}

	claims, err := pkg.VerifyToken(token)
	if err != nil {
		c.JSON(403, "can't connect to websocket because access token is invalid")
		return
	}

	c.Request.Header.Set("X-User-ID", claims.UserID.String())
	wh.mel.HandleRequest(c.Writer, c.Request)

}

func (wh *WebsocketHandler) RegisterRoutes(r gin.IRouter) {
	wsEndpoint := r.Group("/ws")
	wsEndpoint.GET("/", wh.HandleHandshakeRequest)
}
