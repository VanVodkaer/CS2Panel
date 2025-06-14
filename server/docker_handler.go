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
		// 容器名称（必填）
		Name string `json:"name" binding:"required"`

		// 校验参数（可选，默认值为 "1"）
		STEAMAPPVALIDATE string `json:"steamappvalidate"` // "0" 不校验，"1" 校验

		// 以下参数均为可选，若未提供则使用默认值
		CS2_PORT             string `json:"cs2_port"`             // 游戏服务器端口，默认值 "27015"
		CS2_RCON_PORT        string `json:"cs2_rcon_port"`        // RCON 端口，默认值 "27015"
		TV_PORT              string `json:"tv_port"`              // SourceTV 端口，默认值 "27020"
		CS2_LAN              string `json:"cs2_lan"`              // "0" 关闭，"1" 开启局域网模式，默认值 "0"
		CS2_MAXPLAYERS       string `json:"cs2_maxplayers"`       // 最大玩家数
		CS2_STARTMAP         string `json:"cs2_startmap"`         // 启动地图，例如 "de_inferno"
		CS2_MAPGROUP         string `json:"cs2_mapgroup"`         // 地图组名称，例如 "mg_active"
		CS2_SERVERNAME       string `json:"cs2_servername"`       // 服务器名称
		CS2_RCONPW           string `json:"cs2_rconpw"`           // RCON 密码
		CS2_PW               string `json:"cs2_pw"`               // 服务器连接密码
		CS2_CHEATS           string `json:"cs2_cheats"`           // "0" 禁止作弊，"1" 允许作弊，默认值 "0"
		CS2_TV_ENABLE        string `json:"cs2_tv_enable"`        // "0" 禁用，"1" 启用 SourceTV，默认值 "0"
		CS2_TV_PW            string `json:"cs2_tv_pw"`            // SourceTV 观看密码
		CS2_TV_DELAY         string `json:"cs2_tv_delay"`         // SourceTV 延迟，单位为秒
		CS2_TV_AUTORECORD    string `json:"cs2_tv_autorecord"`    // "0" 禁用，"1" 启用 SourceTV 自动录制，默认值 "0"
		CS2_BOT_QUOTA        string `json:"cs2_bot_quota"`        // 机器人数量
		CS2_BOT_DIFFICULTY   string `json:"cs2_bot_difficulty"`   // "0" 最容易，"3" 最难，默认值 "1"
		CS2_COMPETITIVE_MODE string `json:"cs2_competitive_mode"` // "0" 启用，"1" 禁用比赛模式，默认值 "0"
		CS2_LOGGING_ENABLED  string `json:"cs2_logging_enabled"`  // "0" 禁用，"1" 启用日志记录，默认值 "1"
		CS2_GAMEMODE         string `json:"cs2_gamemode"`         // "0" 休闲模式，"1" 竞技模式，默认值 "0"
		CS2_GAMETYPE         string `json:"cs2_gametype"`         // "0" 普通游戏，"1" 死亡竞赛，默认值 "0"
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
			// 校验参数
			fmt.Sprintf("STEAMAPPVALIDATE=%s", util.DefaultIfEmpty(req.STEAMAPPVALIDATE, "0")),
			// 固定环境变量
			fmt.Sprintf("SRCDS_TOKEN=%s", config.GlobalConfig.Game.SRCDS_TOKEN),
			fmt.Sprintf("CS2_PORT=%s", util.DefaultIfEmpty(req.CS2_PORT, "27015")),
			fmt.Sprintf("CS2_RCON_PORT=%s", util.DefaultIfEmpty(req.CS2_RCON_PORT, "27015")),
			fmt.Sprintf("TV_PORT=%s", util.DefaultIfEmpty(req.TV_PORT, "27020")),
			fmt.Sprintf("CS2_SERVERNAME=%s", util.DefaultIfEmpty(req.CS2_SERVERNAME, "Van_Vodkaer's CS2 Server")),
			fmt.Sprintf("CS2_PW=%s", util.DefaultIfEmpty(req.CS2_PW, "")),
			fmt.Sprintf("CS2_RCONPW=%s", util.DefaultIfEmpty(req.CS2_RCONPW, config.GlobalConfig.Game.RCON_PASSWORD)),
			fmt.Sprintf("CS2_LAN=%s", util.DefaultIfEmpty(req.CS2_LAN, "0")),
			fmt.Sprintf("CS2_MAXPLAYERS=%s", util.DefaultIfEmpty(req.CS2_MAXPLAYERS, "")),
			fmt.Sprintf("CS2_STARTMAP=%s", util.DefaultIfEmpty(req.CS2_STARTMAP, "de_dust2")),
			fmt.Sprintf("CS2_MAPGROUP=%s", util.DefaultIfEmpty(req.CS2_MAPGROUP, "mg_active")),
			fmt.Sprintf("CS2_CHEATS=%s", util.DefaultIfEmpty(req.CS2_CHEATS, "0")),
			fmt.Sprintf("CS2_TV_ENABLE=%s", util.DefaultIfEmpty(req.CS2_TV_ENABLE, "0")),
			fmt.Sprintf("CS2_TV_PW=%s", util.DefaultIfEmpty(req.CS2_TV_PW, "")),
			fmt.Sprintf("CS2_TV_DELAY=%s", util.DefaultIfEmpty(req.CS2_TV_DELAY, "0")),
			fmt.Sprintf("CS2_TV_AUTORECORD=%s", util.DefaultIfEmpty(req.CS2_TV_AUTORECORD, "0")),
			fmt.Sprintf("CS2_BOT_QUOTA=%s", util.DefaultIfEmpty(req.CS2_BOT_QUOTA, "0")),
			fmt.Sprintf("CS2_BOT_DIFFICULTY=%s", util.DefaultIfEmpty(req.CS2_BOT_DIFFICULTY, "1")),
			fmt.Sprintf("CS2_COMPETITIVE_MODE=%s", util.DefaultIfEmpty(req.CS2_COMPETITIVE_MODE, "0")),
			fmt.Sprintf("CS2_LOGGING_ENABLED=%s", util.DefaultIfEmpty(req.CS2_LOGGING_ENABLED, "1")),
			fmt.Sprintf("CS2_GAMEMODE=%s", util.DefaultIfEmpty(req.CS2_GAMEMODE, "0")),
			fmt.Sprintf("CS2_GAMETYPE=%s", util.DefaultIfEmpty(req.CS2_GAMETYPE, "0")),
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

// dockerContainerStartHandler 处理启动一个或多个 Docker 容器的请求，并可选地执行命令
func dockerContainerStartHandler(c *gin.Context) {
	// 请求参数结构体：name 或 names 至少提供其一；cmds 可选
	type ContainerStartRequest struct {
		Name  string   `json:"name"`
		Names []string `json:"names"`
		Cmds  []string `json:"cmds"`
	}

	var req ContainerStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}

	// 收集待处理的容器名称
	var targets []string
	if len(req.Names) > 0 {
		targets = req.Names
	} else if req.Name != "" {
		targets = []string{req.Name}
	} else {
		handleErrorResponse(c, "必须提供 name 或 names 参数", nil)
		return
	}

	// 启动成功的容器列表，以及每个容器的命令执行结果
	var started []string
	results := make(map[string][]string)

	for _, name := range targets {
		fullName := FullName(name)
		// 启动容器
		if err := docker.Cli.ContainerStart(context.Background(), fullName, container.StartOptions{}); err != nil {
			handleErrorResponse(c, fmt.Sprintf("启动容器 %s 失败", name), err)
			return
		}
		util.Info(fmt.Sprintf("容器启动成功 容器 ID: %s", name))
		started = append(started, name)

		// 如果传入了 cmds，则对该容器执行命令
		if len(req.Cmds) > 0 {
			responses, err := ExecRconCommands(fullName, req.Cmds)
			if err != nil {
				handleErrorResponse(c, fmt.Sprintf("在容器 %s 中执行命令失败", name), err)
				return
			}
			util.Info(fmt.Sprintf("执行命令成功 容器: %s 命令: %v 响应: %v", name, req.Cmds, responses))
			results[name] = responses
		}
	}

	// 根据是否执行命令，返回不同结构的 JSON
	if len(req.Cmds) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"message":   "容器启动并执行命令成功",
			"started":   started,
			"responses": results,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "容器启动成功",
			"started": started,
		})
	}
}

// dockerContainerStopHandler 处理停止一个或多个 Docker 容器的请求
func dockerContainerStopHandler(c *gin.Context) {
	// 请求参数结构体：name 或 names 至少提供其一
	type ContainerStopRequest struct {
		Name  string   `json:"name"`
		Names []string `json:"names"`
	}

	var req ContainerStopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}

	// 收集待停止的容器名称
	var targets []string
	if len(req.Names) > 0 {
		targets = req.Names
	} else if req.Name != "" {
		targets = []string{req.Name}
	} else {
		handleErrorResponse(c, "必须提供 name 或 names 参数", nil)
		return
	}

	// 存放已成功停止的容器列表
	var stopped []string

	for _, name := range targets {
		fullName := FullName(name)
		// 停止容器
		if err := docker.Cli.ContainerStop(context.Background(), fullName, container.StopOptions{}); err != nil {
			handleErrorResponse(c, fmt.Sprintf("停止容器 %s 失败", name), err)
			return
		}
		util.Info(fmt.Sprintf("容器停止成功 容器 ID: %s", name))
		stopped = append(stopped, name)
	}

	// 返回停止成功的消息和列表
	c.JSON(http.StatusOK, gin.H{
		"message": "容器停止成功",
		"stopped": stopped,
	})
}

// dockerContainerRestartHandler 处理重启一个或多个 Docker 容器的请求
func dockerContainerRestartHandler(c *gin.Context) {
	// 请求参数结构体：name 或 names 至少提供其一
	type ContainerRestartRequest struct {
		Name  string   `json:"name"`
		Names []string `json:"names"`
	}

	// 解析请求体
	var req ContainerRestartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}

	// 收集待重启的容器名称
	var targets []string
	if len(req.Names) > 0 {
		targets = req.Names
	} else if req.Name != "" {
		targets = []string{req.Name}
	} else {
		handleErrorResponse(c, "必须提供 name 或 names 参数", nil)
		return
	}

	// 存放已成功重启的容器列表
	var restarted []string

	for _, name := range targets {
		fullName := FullName(name)
		// 重启容器（如需传超时时间可自行拓展 RestartOptions）
		if err := docker.Cli.ContainerRestart(context.Background(), fullName, container.StopOptions{}); err != nil {
			handleErrorResponse(c, fmt.Sprintf("重启容器 %s 失败", name), err)
			return
		}
		util.Info(fmt.Sprintf("容器重启成功 容器 ID: %s", name))
		restarted = append(restarted, name)
	}

	// 返回重启成功的消息和列表
	c.JSON(http.StatusOK, gin.H{
		"message":   "容器重启成功",
		"restarted": restarted,
	})
}

// dockerContainerRemoveHandler 处理删除 Docker 容器的请求
func dockerContainerRemoveHandler(c *gin.Context) {
	// 定义请求参数结构体
	type ContainerRemoveRequest struct {
		Name  string   `json:"name"`
		Names []string `json:"names"`
	}

	// 从请求中解析参数
	var req ContainerRemoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}

	// 收集待删除的容器名称
	var targets []string
	if len(req.Names) > 0 {
		targets = req.Names
	} else if req.Name != "" {
		targets = []string{req.Name}
	} else {
		handleErrorResponse(c, "必须提供 name 或 names 参数", nil)
		return
	}

	// 存放已成功删除的容器列表
	var removed []string

	for _, name := range targets {
		fullName := FullName(name)
		// 先停止容器
		if err := docker.Cli.ContainerStop(context.Background(), fullName, container.StopOptions{}); err != nil {
			handleErrorResponse(c, fmt.Sprintf("停止容器 %s 失败", name), err)
			return
		}
		// 删除容器
		if err := docker.Cli.ContainerRemove(context.Background(), fullName, container.RemoveOptions{}); err != nil {
			handleErrorResponse(c, fmt.Sprintf("删除容器 %s 失败", name), err)
			return
		}
		util.Info(fmt.Sprintf("容器删除成功 容器 ID: %s", name))
		removed = append(removed, name)
	}

	// 返回删除成功的消息和列表
	c.JSON(http.StatusOK, gin.H{
		"message": "容器删除成功",
		"removed": removed,
	})

	util.Info("容器删除成功 容器 ID: " + req.Name)
}
