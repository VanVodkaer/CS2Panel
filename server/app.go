package server

import (
	"fmt"

	"github.com/VanVodkaer/CS2Panel/config"
	"github.com/VanVodkaer/CS2Panel/util"
	"github.com/gin-gonic/gin"
)

// App 结构体 包含 Gin 引擎和应用配置
type App struct {
	Router *gin.Engine
	Config *config.Config
}

func init() { // 将 Gin 设置为 release 模式
	if config.GlobalConfig.Env.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
}

// ServerNewApp 创建并初始化一个新的 Web 应用
func ServerNewApp() (*App, error) {
	// 初始化 Gin 路由器
	router := gin.Default()
	// 设置 Web 应用的路由
	ServerSetRouter(router)

	// 返回初始化的应用
	return &App{
		Router: router,
		Config: config.GlobalConfig,
	}, nil
}

// ServerStart 启动 Web 服务器
func (app *App) ServerStart() {
	// 使用配置中的端口号
	port := fmt.Sprintf(":%d", app.Config.Server.Port)
	util.Info("服务器运行端口" + port)

	// 启动服务器
	if err := app.Router.Run(port); err != nil {
		util.Error("服务器启动失败: %v", err)
	}
}
