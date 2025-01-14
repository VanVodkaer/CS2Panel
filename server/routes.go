package server

import (
	"github.com/gin-gonic/gin"
)

// ServerSetRouter 设置 Web 应用的路由
func ServerSetRouter(router *gin.Engine) {

	api := router.Group("/api")
	{
		docker := api.Group("/docker")
		{
			docker.Any("/ping", pingHandler)
			image := docker.Group("/image")
			{
				image.POST("/pull", imagePullHandler)
			}
			container := docker.Group("/container")
			{
				container.GET("/list", containerListHandler)
				container.POST("/create", containerCreateHandler)
				container.POST("/start", containerStartHandler)
				container.POST("/stop", containerStopHandler)
				container.POST("/restart", containerRestartHandler)
				container.POST("/remove", containerRemoveHandler)
				container.POST("/exec", containerExecHandler)
			}

		}
	}
}
