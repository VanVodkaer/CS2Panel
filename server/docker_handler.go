package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/VanVodkaer/CS2Panel/config"
	"github.com/VanVodkaer/CS2Panel/docker"
	"github.com/VanVodkaer/CS2Panel/util"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/go-connections/nat"
	"github.com/gin-gonic/gin"
)

// pingHandler 处理 Docker 服务的 ping 请求
func pingHandler(c *gin.Context) {
	ping, err := docker.Cli.Ping(context.Background())
	if err != nil {
		handleErrorResponse(c, "Docker 服务连接失败", err)
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Docker 服务正在运行",
			"ping":    ping,
		})
	}

	util.Info(fmt.Sprintf("Docker 服务正在运行, Ping 信息: %+v", ping))
	return
}

// imagePullHandler 处理拉取 Docker 镜像的请求
func imagePullHandler(c *gin.Context) {
	reader, err := docker.Cli.ImagePull(context.Background(), config.GlobalConfig.Docker.ImageName, image.PullOptions{})
	if err != nil {
		handleErrorResponse(c, "拉取 Docker 镜像失败", err)
		return
	}
	defer reader.Close()

	c.JSON(http.StatusOK, gin.H{
		"message": "拉取 Docker 镜像成功",
	})

	util.Info("拉取 Docker 镜像成功")
	return
}

// containerListHandler 处理获取 Docker 容器列表的请求
func containerListHandler(c *gin.Context) {
	// 过滤器按容器名称前缀过滤
	containers, err := docker.Cli.ContainerList(context.Background(), container.ListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("name", config.GlobalConfig.Docker.Prefix)),
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

	util.Info("获取 Docker 容器列表成功")
	return
}

// DockerCreateRequest 定义创建 Docker 容器的请求参数
type ContainerCreateRequest struct {
	Name      string `json:"name" binding:"required"` // 容器名称
	GamePort  string `json:"game_port"`               // 用于游戏服务器的端口
	WatchPort string `json:"watch_port"`              // 用于观战服务器状态的端口
}

// containerCreateHandler 处理创建 Docker 容器的请求
func containerCreateHandler(c *gin.Context) {
	// 从请求中解析参数
	var req ContainerCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	// 获取当前工作目录
	cwd, err := os.Getwd()
	// 绑定路径拼接
	bindPath := fmt.Sprintf("%s:/home/steam/cs2-dedicated", filepath.Join(cwd, "cs2-data"))

	if err != nil {
		handleErrorResponse(c, "获取当前工作目录失败", err)
		return
	}
	// 定义容器的创建配置
	containerConfig := &container.Config{
		Image: config.GlobalConfig.Docker.ImageName,
		ExposedPorts: nat.PortSet{
			"27015/tcp": {},
			"27020/tcp": {},
			"27015/udp": {},
			"27020/udp": {},
		},
		Env: []string{
			fmt.Sprintf("SRCDS_TOKEN=%s", os.Getenv("SRCDS_TOKEN")),
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"27015/tcp": []nat.PortBinding{{
				HostPort: util.DefaultIfEmpty(req.GamePort, "27015"),
			}},
			"27020/tcp": []nat.PortBinding{{
				HostPort: util.DefaultIfEmpty(req.WatchPort, "27020"),
			}},
			"27015/udp": []nat.PortBinding{{
				HostPort: util.DefaultIfEmpty(req.GamePort, "27015"),
			}},
			"27020/udp": []nat.PortBinding{{
				HostPort: util.DefaultIfEmpty(req.WatchPort, "27020"),
			}},
		},
		Binds: []string{
			bindPath,
		},
	}

	// 创建容器
	createResp, err := docker.Cli.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, nil, config.GlobalConfig.Docker.Prefix+"-"+req.Name)
	if err != nil {
		handleErrorResponse(c, "创建容器失败", err)
		return
	} else {
		// 返回容器创建成功的消息和容器 ID
		c.JSON(200, gin.H{
			"message":      "容器创建成功",
			"container_id": createResp.ID,
		})
	}

	util.Info("容器创建成功 容器 ID: " + createResp.ID)
	return
}

// ContainerStartRequest 定义启动 Docker 容器的请求参数
type ContainerStartRequest struct {
	ID string `json:"id" binding:"required"`
}

// containerStartHandler 处理启动 Docker 容器的请求
func containerStartHandler(c *gin.Context) {
	// 从请求中解析参数
	var req ContainerStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	// 启动容器
	if err := docker.Cli.ContainerStart(context.Background(), config.GlobalConfig.Docker.Prefix+"-"+req.ID, container.StartOptions{}); err != nil {
		handleErrorResponse(c, "启动容器失败", err)
		return
	} else {
		// 返回容器启动成功的消息
		c.JSON(200, gin.H{
			"message": "容器启动成功",
		})
	}

	util.Info("容器启动成功 容器 ID: " + req.ID)
	return
}

// ContainerStopRequest 定义停止 Docker 容器的请求参数
type ContainerStopRequest struct {
	ID string `json:"id" binding:"required"`
}

// containerStopHandler 处理停止 Docker 容器的请求
func containerStopHandler(c *gin.Context) {
	// 从请求中解析参数
	var req ContainerStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	// 停止容器
	if err := docker.Cli.ContainerStop(context.Background(), req.ID, container.StopOptions{}); err != nil {
		handleErrorResponse(c, "停止容器失败", err)
		return
	} else {
		// 返回容器停止成功的消息
		c.JSON(200, gin.H{
			"message": "容器停止成功",
		})
	}

	util.Info("容器停止成功 容器 ID: " + req.ID)
	return
}

// ContainerRestartRequest 定义重启 Docker 容器的请求参数
type ContainerRestartRequest struct {
	ID string `json:"id" binding:"required"`
}

// containerRestartHandler 处理重启 Docker 容器的请求
func containerRestartHandler(c *gin.Context) {
	// 从请求中解析参数
	var req ContainerStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	// 重启容器
	if err := docker.Cli.ContainerRestart(context.Background(), req.ID, container.StopOptions{}); err != nil {
		handleErrorResponse(c, "重启容器失败", err)
		return
	} else {
		// 返回容器重启成功的消息
		c.JSON(200, gin.H{
			"message": "容器重启成功",
		})
	}

	util.Info("容器重启成功 容器 ID: " + req.ID)
	return
}

// ContainerRemoveRequest 定义删除 Docker 容器的请求参数
type ContainerRemoveRequest struct {
	ID string `json:"id" binding:"required"`
}

// containerRemoveHandler 处理删除 Docker 容器的请求
func containerRemoveHandler(c *gin.Context) {
	// 从请求中解析参数
	var req ContainerStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	// 删除容器
	if err := docker.Cli.ContainerRemove(context.Background(), req.ID, container.RemoveOptions{}); err != nil {
		handleErrorResponse(c, "删除容器失败", err)
		return
	} else {
		// 返回容器删除成功的消息
		c.JSON(200, gin.H{
			"message": "容器删除成功",
		})
	}

	util.Info("容器删除成功 容器 ID: " + req.ID)
	return
}

// ContainerExecRequest 定义执行 Docker 容器命令的请求参数
type ContainerExecRequest struct {
	ID  string   `json:"id" binding:"required"`
	Cmd []string `json:"cmd" binding:"required"`
}

// containerExecHandler 处理执行 Docker 容器命令的请求
func containerExecHandler(c *gin.Context) {
	// 从请求中解析参数
	var req ContainerExecRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	// 创建执行命令的配置
	execConfig := container.ExecOptions{
		Cmd: req.Cmd,
	}
	// 执行命令
	execID, err := docker.Cli.ContainerExecCreate(context.Background(), req.ID, execConfig)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		// 返回执行命令成功的消息和执行 ID
		c.JSON(200, gin.H{
			"message": "执行命令成功",
			"exec_id": execID.ID,
		})
	}

	util.Info("执行命令成功 容器 ID: " + req.ID + " 命令: " + fmt.Sprintf("%v", req.Cmd))
	return

}
