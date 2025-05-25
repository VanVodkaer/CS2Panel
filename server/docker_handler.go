package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

// dockerPingHandler 处理 Docker 服务的 ping 请求
func dockerPingHandler(c *gin.Context) {
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

// dockerImagePullHandler 异步处理拉取 Docker 镜像的请求
var pullStatus sync.Map // key: imageName, value: status(string)
func dockerImagePullHandler(c *gin.Context) {
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

// dockerImagePullStatusHandler 处理获取 Docker 镜像拉取状态的请求
func dockerImagePullStatusHandler(c *gin.Context) {
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

// dockerContainerListHandler 处理获取 Docker 容器列表的请求
func dockerContainerListHandler(c *gin.Context) {
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

// dockerContainerCreateHandler 处理创建 Docker 容器的请求
func dockerContainerCreateHandler(c *gin.Context) {
	// 定义请求参数结构体
	type ContainerCreateRequest struct {
		// 容器名称 (必填)
		Name string `json:"name" binding:"required"`

		// 容器创建时确定的参数
		// 用于游戏服务器和Rcon的端口 (可选，默认值为 "27015")
		CS2_PORT string `json:"cs2_port"` // 游戏服务器端口，默认值 "27015"
		// 用于Rcon的端口 (可选，默认值为 "27015")
		CS2_RCON_PORT string `json:"cs2_rcon_port"` // RCON 端口，默认值 "27015"
		// 用于观战服务器状态的端口 (可选，默认值为 "27020")
		TV_PORT string `json:"tv_port"` // SourceTV 端口，默认值 "27020"
		// 是否为局域网模式 (可选，默认值为 "0"，"0"为局域网模式，"1"为非局域网模式)
		CS2_LAN string `json:"cs2_lan"` // 局域网模式，默认值 "0"
		// 最大玩家数 (可选)
		CS2_MAXPLAYERS string `json:"cs2_maxplayers"` // 最大玩家数
		// 游戏开始地图 (可选)
		CS2_STARTMAP string `json:"cs2_startmap"` // 启动地图，例如 "de_inferno"
		// 地图组 (可选)
		CS2_MAPGROUP string `json:"cs2_mapgroup"` // 地图组名称，例如 "mg_active"

		// 容器运行时的参数（这些可以通过控制台或 RCON 动态修改）
		// 服务器名称
		CS2_SERVERNAME string `json:"cs2_servername"` // 服务器名称
		// RCON 密码
		CS2_RCONPW string `json:"cs2_rconpw"` // RCON 密码
		// 连接密码
		CS2_PW string `json:"cs2_pw"` // 服务器连接密码
		// 是否允许作弊 (可选，默认值为 "0"，"0" 禁止作弊，"1" 允许作弊)
		CS2_CHEATS string `json:"cs2_cheats"` // 作弊模式，默认值 "0"
		// 是否启用 SourceTV (可选，默认值为 "0"，"0" 禁用，"1" 启用)
		CS2_TV_ENABLE string `json:"cs2_tv_enable"` // 启用 SourceTV，默认值 "0"
		// SourceTV 密码
		CS2_TV_PW string `json:"cs2_tv_pw"` // SourceTV 观看密码
		// SourceTV 延迟 (可选)
		CS2_TV_DELAY string `json:"cs2_tv_delay"` // SourceTV 延迟，单位为秒
		// 自动录制 SourceTV (可选，默认值为 "0"，"0" 禁用，"1" 启用)
		CS2_TV_AUTORECORD string `json:"cs2_tv_autorecord"` // 启用 SourceTV 自动录制，默认值 "0"
		// 机器人数量 (可选，默认值为 "0")
		CS2_BOT_QUOTA string `json:"cs2_bot_quota"` // 机器人数量
		// 机器人难度 (可选，默认值为 "1"，"0" 最容易，"3" 最难)
		CS2_BOT_DIFFICULTY string `json:"cs2_bot_difficulty"` // 机器人难度，默认值 "1"
		// 是否开启比赛模式 (可选，默认值为 "0" 启用，"1" 禁用)
		CS2_COMPETITIVE_MODE string `json:"cs2_competitive_mode"` // 比赛模式，默认值 "0"
		// 日志记录是否启用 (可选，默认值为 "1"，"1" 启用，"0" 禁用)
		CS2_LOGGING_ENABLED string `json:"cs2_logging_enabled"` // 日志记录启用，默认值 "1"
		// 游戏模式 (可选，默认值为 "0"，"0" 为休闲模式，"1" 为竞技模式)
		CS2_GAMEMODE string `json:"cs2_gamemode"` // 游戏模式，默认值 "0"
		// 游戏类型 (可选，默认值为 "0"，"0" 为普通游戏，"1" 为死亡竞赛)
		CS2_GAMETYPE string `json:"cs2_gametype"` // 游戏类型，默认值 "0"
	}

	// 从请求中解析参数
	var req ContainerCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}

	// 定义容器的创建配置
	containerConfig := &container.Config{
		Image: config.GlobalConfig.Docker.ImageName,
		ExposedPorts: nat.PortSet{
			nat.Port(fmt.Sprintf("%s/tcp", util.DefaultIfEmpty(req.CS2_RCON_PORT, "27015"))): {},
			nat.Port(fmt.Sprintf("%s/udp", util.DefaultIfEmpty(req.CS2_PORT, "27015"))):      {},
			nat.Port(fmt.Sprintf("%s/udp", util.DefaultIfEmpty(req.TV_PORT, "27020"))):       {},
		},
		Env: []string{
			// 固定环境变量
			fmt.Sprintf("SRCDS_TOKEN=%s", config.GlobalConfig.Game.SRCDS_TOKEN),
			// 端口信息
			fmt.Sprintf("CS2_PORT=%s", util.DefaultIfEmpty(req.CS2_PORT, "27015")),
			fmt.Sprintf("CS2_RCON_PORT=%s", util.DefaultIfEmpty(req.CS2_RCON_PORT, "27015")),
			fmt.Sprintf("TV_PORT=%s", util.DefaultIfEmpty(req.TV_PORT, "27020")),
			// 其它环境变量
			fmt.Sprintf("CS2_SERVERNAME=%s", util.DefaultIfEmpty(req.CS2_SERVERNAME, "Van_Vodkaer's CS2 Server")),
			fmt.Sprintf("CS2_PW=%s", util.DefaultIfEmpty(req.CS2_PW, "")),
			fmt.Sprintf("CS2_RCONPW=%s", util.DefaultIfEmpty(req.CS2_RCONPW, config.GlobalConfig.Game.RCON_PASSWORD)),
			// 其他环境变量的定义
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port(fmt.Sprintf("%s/tcp", util.DefaultIfEmpty(req.CS2_RCON_PORT, "27015"))): []nat.PortBinding{{
				HostPort: util.DefaultIfEmpty(req.CS2_RCON_PORT, "27015"),
			}},
			nat.Port(fmt.Sprintf("%s/udp", util.DefaultIfEmpty(req.CS2_PORT, "27015"))): []nat.PortBinding{{
				HostPort: util.DefaultIfEmpty(req.CS2_PORT, "27015"),
			}},
			nat.Port(fmt.Sprintf("%s/udp", util.DefaultIfEmpty(req.TV_PORT, "27020"))): []nat.PortBinding{{
				HostPort: util.DefaultIfEmpty(req.TV_PORT, "27020"),
			}},
		},
		Binds: []string{
			fmt.Sprintf("%s:/home/steam/cs2-dedicated", config.GlobalConfig.Docker.VolumeName),
		},
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

// dockerContainerStartHandler 处理启动 Docker 容器的请求
func dockerContainerStartHandler(c *gin.Context) {
	// 定义请求参数结构体
	type ContainerStartRequest struct {
		Name string   `json:"name" binding:"required"`
		Cmds []string `json:"cmds"` // 可选的命令参数
	}

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
		responses, err := ExecRconCommands(FullName(req.Name), req.Cmds)
		if err != nil {
			handleErrorResponse(c, "执行命令失败", err)
			return
		} else {
			util.Info("执行命令成功 命令: " + fmt.Sprintf("%v", req.Cmds) + " 响应: " + fmt.Sprintf("%v", responses))
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

// dockerContainerStopHandler 处理停止 Docker 容器的请求
func dockerContainerStopHandler(c *gin.Context) {
	// 定义请求参数结构体
	type ContainerStopRequest struct {
		Name string `json:"name" binding:"required"`
	}

	// 从请求中解析参数
	var req ContainerStopRequest
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

// dockerContainerRestartHandler 处理重启 Docker 容器的请求
func dockerContainerRestartHandler(c *gin.Context) {
	// 定义请求参数结构体
	type ContainerRestartRequest struct {
		Name string `json:"name" binding:"required"`
	}

	// 从请求中解析参数
	var req ContainerRestartRequest
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

// dockerContainerRemoveHandler 处理删除 Docker 容器的请求
func dockerContainerRemoveHandler(c *gin.Context) {
	// 定义请求参数结构体
	type ContainerRemoveRequest struct {
		Name string `json:"name" binding:"required"`
	}

	// 从请求中解析参数
	var req ContainerRemoveRequest
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
