package server

import (
	"github.com/gin-gonic/gin"
)

// ServerSetRouter 设置 Web 应用的路由
func ServerSetRouter(router *gin.Engine) {
	router.GET("/", rootHandler)
}
