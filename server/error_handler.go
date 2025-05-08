package server

import (
	"net/http"

	"github.com/VanVodkaer/CS2Panel/util"
	"github.com/gin-gonic/gin"
)

// handleErrorResponse 处理错误响应，记录错误日志并返回 500 错误给客户端
func handleErrorResponse(c *gin.Context, message string, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error":   message,
		"details": err.Error(), // 将错误详情返回给客户端，便于调试
	})

	util.Error(message, err) // 日志记录
}
