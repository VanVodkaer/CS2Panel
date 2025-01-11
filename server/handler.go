package server

import (
	"context"
	"net/http"

	"github.com/VanVodkaer/CS2Panel/config"
	"github.com/VanVodkaer/CS2Panel/docker"
	"github.com/VanVodkaer/CS2Panel/util"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/gin-gonic/gin"
)

// handleErrorResponse 处理错误响应，记录错误日志并返回 500 错误给客户端
func handleErrorResponse(c *gin.Context, message string, err error) {
	util.Error(message, err) // 日志记录
	c.JSON(http.StatusInternalServerError, gin.H{
		"error":   message,
		"details": err.Error(), // 将错误详情返回给客户端，便于调试
	})
}

// dockerlistHandler 处理获取 Docker 容器列表的请求
func dockerlistHandler(c *gin.Context) {
	// 创建一个过滤器，用于按容器名称前缀过滤
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", config.GlobalConfig.Docker.Prefix)

	containers, err := docker.Cli.ContainerList(context.Background(), container.ListOptions{
		All:     true,
		Filters: filterArgs,
	})
	if err != nil {
		handleErrorResponse(c, "获取 Docker 容器列表失败", err)
		return
	}

	if len(containers) == 0 { // 没有容器，返回一个空的容器列表
		c.JSON(http.StatusOK, gin.H{
			"containers": []string{}, // 空列表
		})
	} else { // 返回成功响应并包含指定前缀容器列表
		c.JSON(http.StatusOK, gin.H{
			"containers": containers,
		})
	}

	return
}

func dockercreateHandler(c *gin.Context) {

}
