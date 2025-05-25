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
			dockerGroup.Any("/ping", dockerPingHandler)

			imageGroup := dockerGroup.Group("/image")
			{
				imageGroup.POST("/pull", dockerImagePullHandler)
				imageGroup.GET("/pull/status", dockerImagePullStatusHandler)
			}

			containerGroup := dockerGroup.Group("/container")
			{
				containerGroup.GET("/list", dockerContainerListHandler)
				containerGroup.POST("/create", dockerContainerCreateHandler)
				containerGroup.POST("/start", dockerContainerStartHandler)
				containerGroup.POST("/stop", dockerContainerStopHandler)
				containerGroup.POST("/restart", dockerContainerRestartHandler)
				containerGroup.DELETE("/remove", dockerContainerRemoveHandler)
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
				networkGroup.GET("/addr", infoNetworkAddrHandler)
				networkGroup.GET("/gameport", infoNetworkGamePortHandler)
				networkGroup.GET("/gameports", infoNetworkGamePortsHandler)
				networkGroup.GET("/tvport", infoNetworkTVPortHandler)
				networkGroup.GET("/tvports", infoNetworkTVPortsHandler)
				networkGroup.GET("/gamepasswd", infoNetworkGamePasswdHandler)
				networkGroup.GET("/gamepasswds", infoNetworkGamePasswdsHandler)
				networkGroup.GET("/tvpasswd", infoNetworkTVPasswdHandler)
				networkGroup.GET("/tvpasswds", infoNetworkTVPasswdsHandler)
			}
		}
		rconGroup := apiGroup.Group("/rcon")
		{
			rconGroup.POST("/exec", rconExecHandler)
			gameGroup := rconGroup.Group("/game")
			{
				gameGroup.GET("/status", rconGameStatusHandler)
				gameGroup.POST("/restart", rconGameRestartHandler)
				warmGroup := gameGroup.Group("/warm")
				{
					warmGroup.POST("/start", rconGameWarmStartHandler)
					warmGroup.POST("/end", rconGameWarmEndHandler)
					warmGroup.POST("/offine", rconGameWarmOffineHandler)
					warmGroup.POST("/time", rconGameWarmTimeHandler)
					warmGroup.POST("/pause", rconGameWarmPauseHandler)
				}
				configGroup := rconGroup.Group("/config")
				{
					configGroup.POST("/maxrounds", rconGameConfigMaxRoundsHandler)
					configGroup.POST("/timelimit", rconGameConfigTimeLimitHandler)
					configGroup.POST("/roundtime", rconGameConfigRoundTimeHandler)
					configGroup.POST("/freezetime", rconGameConfigFreezetimeHandler)
					configGroup.POST("/buytime", rconGameConfigBuytimeHandler)
					configGroup.POST("/buy_anywhere", rconGameConfigBuyAnywhereHandler)
					configGroup.POST("/startmoney", rconGameConfigStartMoneyHandler)
					configGroup.POST("/maxmoney", rconGameConfigMaxMoneyHandler)
					configGroup.POST("/autoteambalance", rconGameConfigAutoTeamBalanceHandler)
					configGroup.POST("/limitteams", rconGameConfigLimitTeamsHandler)
					configGroup.POST("/c4timer", rconGameConfigC4TimerHandler)
				}
			}
			mapGroup := rconGroup.Group("/map")
			{
				mapGroup.GET("/now", rconMapNowHandler)
				mapGroup.POST("/change", rconMapChangeHandler)
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
