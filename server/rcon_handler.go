package server

import (
	"encoding/json"

	"github.com/VanVodkaer/CS2Panel/config"
	"github.com/VanVodkaer/CS2Panel/util"
	"github.com/gin-gonic/gin"
)

// containerExecHandler 执行命令
func rconExecHandler(c *gin.Context) {
	// 定义请求参数结构体
	type ContainerExecRequest struct {
		Name string   `json:"name" binding:"required"`
		Cmds []string `json:"cmds" binding:"required"`
	}

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

// rconGameStatusHandler 获取游戏状态
func rconGameStatusHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameStatusRequest struct {
		Name string `form:"name" binding:"required"`
	}

	var req RconGameStatusRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := GetServerStatus(FullName(req.Name))
	if err != nil {
		handleErrorResponse(c, "获取服务器状态失败", err)
		return
	} else {
		data, err := json.Marshal(response)
		if err != nil {
			util.Error("json.Marshal error: ", err)
		} else {
			if config.GlobalConfig.Env.Mode == "debug" {
				util.Debug("获取游戏状态status成功" + " 响应: " + string(data))
			} else {
				util.Info("获取游戏状态status成功")
			}
		}
	}
	c.JSON(200, gin.H{
		"message": "获取游戏状态成功",
		"status":  response,
	})
}

// rconGameStatusJSONHandler 获取游戏状态 status_json
func rconGameStatusJSONHandler(c *gin.Context) {
	type RconGameStatusJSONRequest struct {
		Name string `form:"name" binding:"required"`
	}
	var req RconGameStatusJSONRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}

	status, err := GetServerStatusJSON(FullName(req.Name))
	if err != nil {
		handleErrorResponse(c, "获取服务器状态失败", err)
		return
	}

	// 记录日志
	if data, err := json.Marshal(status); err == nil {
		if config.GlobalConfig.Env.Mode == "debug" {
			util.Debug("获取游戏状态status_json成功" + " 响应: " + string(data))
		} else {
			util.Info("获取游戏状态status_json成功")
		}
	}

	c.JSON(200, gin.H{
		"message": "获取游戏状态成功",
		"status":  status, // 返回解析后的结构体对象
	})
}

// rconGameRestartHandler 重启游戏
func rconGameRestartHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameRestartRequest struct {
		Name  string `json:"name" binding:"required"`
		Value string `json:"value"`
	}

	var req RconGameRestartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_restartgame "+req.Value)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_restartgame " + req.Value + " 响应: " + response)
	}
	c.JSON(200, gin.H{
		"message":  "执行命令成功",
		"response": response,
	})
}

// rconGameConfigModeHandler
func rconGameConfigModeHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameConfigGameModeRequest struct {
		Name     string `json:"name" binding:"required"`
		GameMode string `json:"gamemode"`
		GameType string `json:"gametype"`
	}

	var req RconGameConfigGameModeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	var responses []string
	if req.GameMode != "" {
		response, err := ExecRconCommand(FullName(req.Name), "game_mode "+req.GameMode)
		if err != nil {
			handleErrorResponse(c, "执行命令失败", err)
			return
		} else {
			util.Info("执行命令成功 命令: game_mode " + req.GameMode + " 响应: " + response)
			responses = append(responses, response)
		}
	}
	if req.GameType != "" {
		response, err := ExecRconCommand(FullName(req.Name), "game_type "+req.GameType)
		if err != nil {
			handleErrorResponse(c, "执行命令失败", err)
			return
		} else {
			util.Info("执行命令成功 命令: game_type " + req.GameType + " 响应: " + response)
			responses = append(responses, response)
		}
	}
	c.JSON(200, gin.H{
		"message":   "执行命令成功",
		"responses": responses,
	})
}

// rconGameWarmStartHandler 立刻切换到热身模式
func rconGameWarmStartHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameWarmStartRequest struct {
		Name string `json:"name" binding:"required"`
	}

	var req RconGameWarmStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_warmup_start")
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_warmup_start 响应: " + response)
	}
	c.JSON(200, gin.H{
		"message":  "执行命令成功",
		"response": response,
	})
}

// rconGameWarmEndHandler  立刻结束热身模式
func rconGameWarmEndHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameWarmEndRequest struct {
		Name string `json:"name" binding:"required"`
	}

	var req RconGameWarmEndRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_warmup_end")
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
	} else {
		util.Info("执行命令成功 命令: warmup_end 响应: " + response)
	}
	c.JSON(200, gin.H{
		"message":  "执行命令成功",
		"response": response,
	})
}

// rconGameWarmTimeHandler 设置热身时间
func rconGameWarmTimeHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameWarmTimeRequest struct {
		Name  string `json:"name" binding:"required"`
		Value string `json:"value"` // 热身时长 无参数返回当前时长
	}

	var req RconGameWarmTimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_warmuptime "+req.Value)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_warmuptime " + req.Value + " 响应: " + response)
	}
	c.JSON(200, gin.H{
		"message":  "执行命令成功",
		"response": response,
	})
}

// rconGameWarmPauseHandler 控制热身时间暂停
func rconGameWarmPauseHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameWarmPauseRequest struct {
		Name  string `json:"name" binding:"required"`
		Value string `json:"value"` // 0/false 关闭; 1/true 开启 无参数返回当前状态
	}

	var req RconGameWarmPauseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_warmup_pausetimer "+req.Value)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_warmup_pausetimer " + req.Value + " 响应: " + response)
	}
	c.JSON(200, gin.H{
		"message":  "执行命令成功",
		"response": response,
	})
}

// rconGameConfigMaxRoundsHandler 设置最大回合数
func rconGameConfigMaxRoundsHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameConfigMaxRoundsRequest struct {
		Name  string `json:"name" binding:"required"`
		Value string `json:"value"` // 最大回合数 无参数返回当前最大回合数
	}

	var req RconGameConfigMaxRoundsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_maxrounds "+req.Value)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_maxrounds " + req.Value + " 响应: " + response)
	}
	c.JSON(200, gin.H{
		"message":  "执行命令成功",
		"response": response,
	})
}

// rconGameConfigTimeLimitHandler 设置比赛时间限制
// 每个游戏的最大持续时间，以分钟为单位。默认情况下，此设置处于禁用状态 (设置为 0)。如果当前地图的总持续时间超过此值，当前地图将结束，下一个地图将开始游戏。
func rconGameConfigTimeLimitHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameConfigTimeLimitRequest struct {
		Name  string `json:"name" binding:"required"`
		Value string `json:"value"` // 比赛时间限制 无参数返回当前时间限制
	}

	var req RconGameConfigTimeLimitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_timelimit "+req.Value)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_timelimit " + req.Value + " 响应: " + response)
	}
	c.JSON(200, gin.H{
		"message":  "执行命令成功",
		"response": response,
	})
}

// rconGameConfigRoundTimeHandler 设置每回合时间
func rconGameConfigRoundTimeHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameConfigRoundTimeRequest struct {
		Name  string `json:"name" binding:"required"`
		Value string `json:"value"` // 回合时间 无参数返回当前回合时间
		Mode  string `json:"mode"`  // 模式 可选参数 defuse 拆弹模式, hostage 人质解救
	}

	var req RconGameConfigRoundTimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	var command string
	if req.Mode == "defuse" {
		command = "mp_roundtime_defuse " + req.Value
	} else if req.Mode == "hostage" {
		command = "mp_roundtime_hostage " + req.Value
	} else {
		command = "mp_roundtime " + req.Value
	}

	response, err := ExecRconCommand(FullName(req.Name), command)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: " + command + " 响应: " + response)
	}
	c.JSON(200, gin.H{
		"message":  "执行命令成功",
		"response": response,
	})
}

// rconGameConfigFreezetimeHandler 设置冻结时间
func rconGameConfigFreezetimeHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameConfigFreezetimeRequest struct {
		Name  string `json:"name" binding:"required"`
		Value string `json:"value"` // 冻结时间 无参数返回当前冻结时间
	}

	var req RconGameConfigFreezetimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_freezetime "+req.Value)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_freezetime " + req.Value + " 响应: " + response)
	}
	c.JSON(200, gin.H{
		"message":  "执行命令成功",
		"response": response,
	})
}

// rconGameConfigBuytimeHandler 设置购买时间
func rconGameConfigBuytimeHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameConfigBuytimeRequest struct {
		Name  string `json:"name" binding:"required"`
		Value string `json:"value"`
	}

	var req RconGameConfigBuytimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_buytime "+req.Value)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_buytime " + req.Value + " 响应: " + response)
	}
	c.JSON(200, gin.H{
		"message":  "执行命令成功",
		"response": response,
	})
}

// rconGameConfigBuyAnywhereHandler 设置是否允许在地图任意位置购买装备
func rconGameConfigBuyAnywhereHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameConfigBuyAnywhereRequest struct {
		Name  string `json:"name" binding:"required"`
		Value string `json:"value"`
	}

	var req RconGameConfigBuyAnywhereRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_buy_anywhere "+req.Value)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_buy_anywhere " + req.Value + " 响应: " + response)
	}
	c.JSON(200, gin.H{
		"message":  "执行命令成功",
		"response": response,
	})
}

// rconGameConfigStartMoneyHandler 设置初始金钱
func rconGameConfigStartMoneyHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameConfigStartMoneyRequest struct {
		Name  string `json:"name" binding:"required"`
		Value string `json:"value"` // 初始金钱 无参数返回当前初始金钱
	}

	var req RconGameConfigStartMoneyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_startmoney "+req.Value)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_startmoney " + req.Value + " 响应: " + response)
	}
	c.JSON(200, gin.H{
		"message":  "执行命令成功",
		"response": response,
	})
}

// rconGameConfigMaxMoneyHandler 设置最大金钱
func rconGameConfigMaxMoneyHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameConfigMaxMoneyRequest struct {
		Name  string `json:"name" binding:"required"`
		Value string `json:"value"` // 最大金钱 无参数返回当前最大金钱
	}

	var req RconGameConfigMaxMoneyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_maxmoney "+req.Value)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_maxmoney " + req.Value + " 响应: " + response)
	}
	c.JSON(200, gin.H{
		"message":  "执行命令成功",
		"response": response,
	})
}

// rconGameConfigAutoTeamBalanceHandler 设置自动队伍平衡
func rconGameConfigAutoTeamBalanceHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameConfigAutoTeamBalanceRequest struct {
		Name  string `json:"name" binding:"required"`
		Value string `json:"value"`
	}

	var req RconGameConfigAutoTeamBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_autoteambalance "+req.Value)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_autoteambalance " + req.Value + " 响应: " + response)
	}
	c.JSON(200, gin.H{
		"message":  "执行命令成功",
		"response": response,
	})
}

// rconGameConfigAutoKickHandler 设置自动踢出空闲玩家
func rconGameConfigAutoKickHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameConfigAutoKickRequest struct {
		Name  string `json:"name" binding:"required"`
		Value string `json:"value"`
	}

	var req RconGameConfigAutoKickRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_autokick "+req.Value)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_autokick " + req.Value + " 响应: " + response)
	}
	c.JSON(200, gin.H{
		"message":  "执行命令成功",
		"response": response,
	})
}

// rconGameConfigLimitTeamsHandler 设置两个队伍之间允许存在的玩家差异数量的最大值，0为无限制
func rconGameConfigLimitTeamsHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameConfigLimitTeamsRequest struct {
		Name  string `json:"name" binding:"required"`
		Value string `json:"value"` // 允许存在的玩家差异数量的最大值 无参数返回当前最大值
	}

	var req RconGameConfigLimitTeamsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_limitteams "+req.Value)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_limitteams " + req.Value + " 响应: " + response)
	}
	c.JSON(200, gin.H{
		"message":  "执行命令成功",
		"response": response,
	})
}

// rconGameConfigC4TimerHandler 设置 C4 爆炸倒计时
func rconGameConfigC4TimerHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameConfigC4TimerRequest struct {
		Name  string `json:"name" binding:"required"`
		Value string `json:"value"` // C4 爆炸倒计时 无参数返回当前倒计时
	}

	var req RconGameConfigC4TimerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_c4timer "+req.Value)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_c4timer " + req.Value + " 响应: " + response)
	}
	c.JSON(200, gin.H{
		"message":  "执行命令成功",
		"response": response,
	})
}

// rconMapNowHandler 获取当前地图
func rconMapNowHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconMapNowRequest struct {
		Name string `form:"name" binding:"required"`
	}

	var req RconMapNowRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := GetServerStatus(FullName(req.Name))
	if err != nil {
		handleErrorResponse(c, "获取服务器状态失败", err)
		return
	} else {
		util.Info("获取当前地图成功 响应: " + response.Spawngroups[0].Path)
	}
	c.JSON(200, gin.H{
		"message": "获取当前地图成功",
		"map":     response.Spawngroups[0].Path,
	})
}

// rconMapChangeHandler
func rconMapChangeHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconMapChangeRequest struct {
		Name string `json:"name" binding:"required"`
		Map  string `json:"map" binding:"required"`
	}

	var req RconMapChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "map "+req.Map)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: map " + req.Map + " 响应: " + response)
	}
	c.JSON(200, gin.H{
		"message":  "执行命令成功",
		"response": response,
	})
}

func rconGameUserKickHandler(c *gin.Context) {
	type RconGameUserKickRequest struct {
		Name string `json:"name" binding:"required"`
		User string `json:"user" binding:"required"`
	}
	var req RconGameUserKickRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}

	util.Debug("rconGameUserKickHandler 请求参数: " + "Name: " + req.Name + ", User: " + req.User)

	response, err := ExecRconCommand(FullName(req.Name), "kick \""+req.User+"\"")
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	}
	util.Info("执行命令成功 命令: kick \"" + req.User + "\" 响应: " + response)

	c.JSON(200, gin.H{
		"message":  "执行命令成功",
		"response": response,
	})
}

// rconGameUserBanIDHandler 禁止玩家ID
func rconGameUserBanIDHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameUserBanIDRequest struct {
		Name string `json:"name" binding:"required"`
		Time string `json:"time" binding:"required"`
		ID   string `json:"id" binding:"required"` // steamID64
	}

	var req RconGameUserBanIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "banid "+req.Time+" \""+req.ID+"\"")
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: banid 0 \"" + req.ID + "\" 响应: " + response)
	}
	c.JSON(200, gin.H{
		"message":  "执行命令成功",
		"response": response,
	})
}
