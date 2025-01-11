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
			docker.GET("/list", dockerlistHandler)
			docker.POST("/create", dockercreateHandler)
			// docker.POST("/start", dockerstartHandler)
			// docker.POST("/stop", dockerstopHandler)
			// docker.POST("/restart", dockerrestartHandler)
			// docker.POST("/remove", dockerremoveHandler)
			// docker.POST("/exec", dockerexecHandler)
		}
	}
}
