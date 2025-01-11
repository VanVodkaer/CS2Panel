package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func rootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "欢迎来到 CS2Panel!",
	})
}
