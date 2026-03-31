package bootstrap

import "github.com/gin-gonic/gin"

type BootstrapHandler interface {
	RegisterRoutes(r gin.IRouter)
}

func SetupRoutes(router *gin.Engine, b ...BootstrapHandler) {
	api := router.Group("/api")
	for _, h := range b {
		h.RegisterRoutes(api)
	}

	// todo : Handle public files path !
}
