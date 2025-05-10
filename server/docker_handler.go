package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"

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
}

// imagePullHandler 异步处理拉取 Docker 镜像的请求
var pullStatus sync.Map // key: imageName, value: status(string)
func imagePullHandler(c *gin.Context) {
	imageName := config.GlobalConfig.Docker.ImageName
	_, loaded := pullStatus.LoadOrStore(imageName, "pulling")
	if loaded {
		c.JSON(http.StatusOK, gin.H{"message": "镜像正在拉取中"})
		return
	}

	go func() {
		defer pullStatus.Delete(imageName)

		reader, err := docker.Cli.ImagePull(context.Background(), imageName, image.PullOptions{})
		if err != nil {
			util.Error("拉取失败：", err)
			pullStatus.Store(imageName, "failed")
			return
		}
		defer reader.Close()

		decoder := json.NewDecoder(reader)
		for decoder.More() {
			var msg map[string]interface{}
			if err := decoder.Decode(&msg); err != nil {
				pullStatus.Store(imageName, "failed")
				return
			}
			if msg["error"] != nil {
				pullStatus.Store(imageName, "failed")
				return
			}
		}

		pullStatus.Store(imageName, "success")
		util.Info("镜像拉取成功")
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"message": "已开始拉取镜像",
	})
}

// imagePullStatusHandler 处理获取 Docker 镜像拉取状态的请求
func imagePullStatusHandler(c *gin.Context) {
	imageName := config.GlobalConfig.Docker.ImageName
	if status, ok := pullStatus.Load(imageName); ok {
		c.JSON(http.StatusOK, gin.H{
			"status": status,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "not_started",
		})
	}
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
}

// DockerCreateRequest 定义创建 Docker 容器的请求参数
type ContainerCreateRequest struct {
	Name       string `json:"name" binding:"required"` // 容器名称
	ServerName string `json:"server_name"`             // 游戏服务器名称
	GamePort   string `json:"game_port"`               // 用于游戏服务器和Rcon的端口
	WatchPort  string `json:"watch_port"`              // 用于观战服务器状态的端口
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
	if err != nil {
		handleErrorResponse(c, "获取当前工作目录失败", err)
		return
	}

	// 路径绑定参数
	csBindPath := fmt.Sprintf("%s:/home/steam/cs2-dedicated", filepath.Join(cwd, config.GlobalConfig.Docker.CSDataDir))

	// 自定义的环境变量
	srcds_token := fmt.Sprintf("SRCDS_TOKEN=%s", config.GlobalConfig.Game.SRCDS_TOKEN)
	cs2_server_name := fmt.Sprintf("CS2_SERVER_NAME=%s", req.ServerName)
	cs2_rconpw := fmt.Sprintf("CS2_RCONPW=%s", config.GlobalConfig.Game.RCON_PASSWORD)

	// 定义容器的创建配置
	containerConfig := &container.Config{
		Image: config.GlobalConfig.Docker.ImageName,
		ExposedPorts: nat.PortSet{
			"27015/tcp": {}, // Rcon端口
			"27015/udp": {}, // 游戏服务器端口
			"27020/udp": {}, // 观战服务器端口
		},
		Env: []string{
			srcds_token,
			cs2_rconpw,
			cs2_server_name,
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"27015/tcp": []nat.PortBinding{{
				HostPort: util.DefaultIfEmpty(req.GamePort, "27015"),
			}},
			"27015/udp": []nat.PortBinding{{
				HostPort: util.DefaultIfEmpty(req.GamePort, "27015"),
			}},
		},
		Binds: []string{
			csBindPath,
		},
	}
	// 如果 WatchPort 不为空，则添加 27020/udp
	if req.WatchPort != "" {
		hostConfig.PortBindings["27020/udp"] = []nat.PortBinding{{
			HostPort: req.WatchPort,
		}}
	}

	// 创建容器
	createResp, err := docker.Cli.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, nil, FullName(req.Name))
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
}

// ContainerStartRequest 定义启动 Docker 容器的请求参数
type ContainerStartRequest struct {
	Name string   `json:"name" binding:"required"`
	Cmds []string `json:"cmds"` // 可选的命令参数
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
	err := docker.Cli.ContainerStart(context.Background(), FullName(req.Name), container.StartOptions{})
	if err != nil {
		handleErrorResponse(c, "启动容器失败", err)
		return
	}
	// 如果提供了命令参数，则执行命令
	if req.Cmds != nil {
		var responses []string

		for _, cmd := range req.Cmds {
			response, err := ExecRconCommand(FullName(req.Name), cmd)
			if err != nil {
				handleErrorResponse(c, "执行命令失败", err)
				return
			} else {
				responses = append(responses, response)
				util.Info("执行命令成功 命令: " + cmd + " 响应: " + response)
			}
		}
		// 返回执行命令的响应
		c.JSON(200, gin.H{
			"message":   "执行命令成功",
			"responses": responses,
		})

	}

	// 返回容器启动成功的消息
	c.JSON(200, gin.H{
		"message": "容器启动成功",
	})

	util.Info("容器启动成功 容器 ID: " + req.Name)
}

// ContainerStopRequest 定义停止 Docker 容器的请求参数
type ContainerStopRequest struct {
	Name string `json:"name" binding:"required"`
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
	if err := docker.Cli.ContainerStop(context.Background(), FullName(req.Name), container.StopOptions{}); err != nil {
		handleErrorResponse(c, "停止容器失败", err)
		return
	} else {
		// 返回容器停止成功的消息
		c.JSON(200, gin.H{
			"message": "容器停止成功",
		})
	}

	util.Info("容器停止成功 容器 ID: " + req.Name)
}

// ContainerRestartRequest 定义重启 Docker 容器的请求参数
type ContainerRestartRequest struct {
	Name string `json:"name" binding:"required"`
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
	if err := docker.Cli.ContainerRestart(context.Background(), FullName(req.Name), container.StopOptions{}); err != nil {
		handleErrorResponse(c, "重启容器失败", err)
		return
	} else {
		// 返回容器重启成功的消息
		c.JSON(200, gin.H{
			"message": "容器重启成功",
		})
	}

	util.Info("容器重启成功 容器 ID: " + req.Name)
}

// ContainerRemoveRequest 定义删除 Docker 容器的请求参数
type ContainerRemoveRequest struct {
	Name string `json:"name" binding:"required"`
}

// containerRemoveHandler 处理删除 Docker 容器的请求
func containerRemoveHandler(c *gin.Context) {
	// 从请求中解析参数
	var req ContainerStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	// 先停止容器
	if err := docker.Cli.ContainerStop(context.Background(), FullName(req.Name), container.StopOptions{}); err != nil {
		handleErrorResponse(c, "停止容器失败", err)
		return
	}
	// 删除容器
	if err := docker.Cli.ContainerRemove(context.Background(), FullName(req.Name), container.RemoveOptions{}); err != nil {
		handleErrorResponse(c, "删除容器失败", err)
		return
	} else {
		// 返回容器删除成功的消息
		c.JSON(200, gin.H{
			"message": "容器删除成功",
		})
	}

	util.Info("容器删除成功 容器 ID: " + req.Name)
}

// ContainerExecRequest 定义执行 Docker 容器命令的请求参数
type ContainerExecRequest struct {
	Name string   `json:"name" binding:"required"`
	Cmds []string `json:"cmds" binding:"required"`
}

// containerExecHandler 处理执行 Docker 容器命令的请求
func containerExecHandler(c *gin.Context) {
	// 从请求中解析参数
	var req ContainerExecRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}

	var responses []string

	for _, cmd := range req.Cmds {
		response, err := ExecRconCommand(FullName(req.Name), cmd)
		if err != nil {
			handleErrorResponse(c, "执行命令失败", err)
			return
		} else {
			responses = append(responses, response)
			util.Info("执行命令成功 命令: " + cmd + " 响应: " + response)
		}
	}
	// 返回执行命令的响应
	c.JSON(200, gin.H{
		"message":   "执行命令成功",
		"responses": responses,
	})
}
