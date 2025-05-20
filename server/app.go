package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/VanVodkaer/CS2Panel/config"
	"github.com/VanVodkaer/CS2Panel/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// App 结构体 包含应用配置
type App struct {
	Config *config.Config
}

// 初始化 Gin 模式
func init() {
	if config.GlobalConfig.Env.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
}

// ServerNewApp 创建并初始化一个新的 Web 应用
func ServerNewApp() (*App, error) {
	return &App{
		Config: config.GlobalConfig,
	}, nil
}

// ServerStart 启动 API 和 Web 服务（根据配置判断是否分端口）
func (app *App) ServerStart() {
	cfg := app.Config

	// 启动 API 服务
	go func() {
		router := gin.Default()
		// 允许所有跨域请求
		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},                            // 允许所有域
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"}, // 允许的 HTTP 方法
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,           // 是否允许带 Cookie
			MaxAge:           12 * time.Hour, // 预检请求的有效期
		}))
		// 注册 API 路由
		ServerSetRouter(router)

		// 如果 Web 服务启用且端口一致，注册 Web 静态路由
		if cfg.Server.WebServer && cfg.Server.WebServerPort == cfg.Server.Port {
			WebServerSetRouter(router)
		}

		addr := fmt.Sprintf(":%d", cfg.Server.Port)
		util.Info("API 服务监听地址: " + addr)

		if err := http.ListenAndServe(addr, router); err != nil {
			util.Error("API 服务启动失败: %v", err)
		}
	}()

	// 启动独立 Web 服务（如果启用且端口不同）
	if cfg.Server.WebServer && cfg.Server.WebServerPort != cfg.Server.Port {
		go func() {
			webRouter := gin.Default()
			WebServerSetRouter(webRouter)

			addr := fmt.Sprintf(":%d", cfg.Server.WebServerPort)
			util.Info("Web 服务监听地址: " + addr)

			if err := http.ListenAndServe(addr, webRouter); err != nil {
				util.Error("Web 服务启动失败: %v", err)
			}
		}()
	} else if !cfg.Server.WebServer {
		util.Warn("Web 服务未启用")
	} else {
		util.Info("Web 服务与 API 服务共用同一端口")
	}

	// 阻塞主线程，防止退出
	select {}
}
