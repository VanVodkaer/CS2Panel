package server

import (
	"github.com/gin-gonic/gin"
)

// ServerSetRouter 设置 API 接口路由
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
				containerGroup.DELETE("/remove", containerRemoveHandler)
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

			networkGroup := infoGroup.Group("/network")
			{
				networkGroup.GET("/addr", networkAddrHandler)
				networkGroup.GET("/gameport", networkGamePortHandler)
				networkGroup.GET("/gameports", networkGamePortsHandler)
				networkGroup.GET("/tvport", networkTVPortHandler)
				networkGroup.GET("/tvports", networkTVPortsHandler)
				networkGroup.GET("/gamepasswd", networkGamePasswdHandler)
				networkGroup.GET("/gamepasswds", networkGamePasswdsHandler)
				networkGroup.GET("/tvpasswd", networkTVPasswdHandler)
				networkGroup.GET("/tvpasswds", networkTVPasswdsHandler)
			}
		}

	}
}

// WebServerSetRouter 设置 Web 静态资源路由
func WebServerSetRouter(router *gin.Engine) {
	// 显式提供静态资源路径，避免 /*filepath 路径冲突
	router.Static("/assets", "./dist/assets")

	// 所有其他路径重定向到 index.html，支持前端路由刷新
	router.NoRoute(func(c *gin.Context) {
		c.File("./dist/index.html")
	})
}
