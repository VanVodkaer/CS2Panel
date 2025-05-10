package server

import (
	"github.com/gin-gonic/gin"
)

// ServerSetRouter 设置 Web 应用的路由
func ServerSetRouter(router *gin.Engine) {

	apiGroup := router.Group("/api")
	{
		dockerGroup := apiGroup.Group("/docker")
		{
			dockerGroup.Any("/ping", pingHandler)
			imageGroup := dockerGroup.Group("/image")
			{
				imageGroup.POST("/pull", imagePullHandler)
				imageGroup.GET("/pull/status", imagePullStatusHandler)
			}
			containerGroup := dockerGroup.Group("/container")
			{
				containerGroup.GET("/list", containerListHandler)
				containerGroup.POST("/create", containerCreateHandler)
				containerGroup.POST("/start", containerStartHandler)
				containerGroup.POST("/stop", containerStopHandler)
				containerGroup.POST("/restart", containerRestartHandler)
				containerGroup.POST("/remove", containerRemoveHandler)
				containerGroup.POST("/exec", containerExecHandler)
			}

		}
		infoGroup := apiGroup.Group("/info")
		{
			mapGroup := infoGroup.Group("/map")
			{
				mapGroup.POST("/update", infoMapUpdateHandler)
				mapGroup.GET("/list", infoMapListHandler)

			}

		}
	}
}
