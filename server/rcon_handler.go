package server

import (
	"encoding/json"

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
			util.Info("获取游戏状态成功 响应: " + string(data))
		}
	}
	c.JSON(200, gin.H{
		"message": "获取游戏状态成功",
		"status":  response,
	})
}

// rconGameRestartHandler 重启游戏
func rconGameRestartHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameRestartRequest struct {
		Name  string `json:"name" binding:"required"`
		Delay string `json:"delay"`
	}

	var req RconGameRestartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_restartgame "+req.Delay)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: restart 响应: " + response)
	}
	c.JSON(200, gin.H{
		"message":  "执行命令成功",
		"response": response,
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

// rconGameWarmOfflineHandler 设置私人游戏 / 离线状态下使用机器人时 热身模式是否开启
func rconGameWarmOfflineHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameWarmOfflineRequest struct {
		Name  string `json:"name" binding:"required"`
		Value string `json:"value"` // 0/false 关闭; 1/true 开启 无参数返回当前状态
	}

	var req RconGameWarmOfflineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_warmup_offline_enabled "+req.Value)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_warmup_offline_enabled " + req.Value + " 响应: " + response)
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
		Name string `json:"name" binding:"required"`
		Time string `json:"time"` // 热身时长 无参数返回当前时长
	}

	var req RconGameWarmTimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_warmup_time "+req.Time)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: warmup_time " + req.Time + " 响应: " + response)
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
		Name      string `json:"name" binding:"required"`
		MaxRounds string `json:"maxrounds"`
	}

	var req RconGameConfigMaxRoundsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_maxrounds "+req.MaxRounds)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_maxrounds " + req.MaxRounds + " 响应: " + response)
	}
	c.JSON(200, gin.H{
		"message":  "执行命令成功",
		"response": response,
	})
}

// rconGameConfigTimeLimitHandler 设置比赛时间限制
func rconGameConfigTimeLimitHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameConfigTimeLimitRequest struct {
		Name      string `json:"name" binding:"required"`
		TimeLimit string `json:"timelimit"`
	}

	var req RconGameConfigTimeLimitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_timelimit "+req.TimeLimit)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_timelimit " + req.TimeLimit + " 响应: " + response)
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
		Name      string `json:"name" binding:"required"`
		RoundTime string `json:"roundtime" binding:"required"`
	}

	var req RconGameConfigRoundTimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_roundtime "+req.RoundTime)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_roundtime " + req.RoundTime + " 响应: " + response)
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
		Name       string `json:"name" binding:"required"`
		FreezeTime string `json:"freezetime" binding:"required"`
	}

	var req RconGameConfigFreezetimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_freezetime "+req.FreezeTime)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_freezetime " + req.FreezeTime + " 响应: " + response)
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
		Name    string `json:"name" binding:"required"`
		BuyTime string `json:"buytime" binding:"required"`
	}

	var req RconGameConfigBuytimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_buytime "+req.BuyTime)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_buytime " + req.BuyTime + " 响应: " + response)
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
		Value string `json:"buy_anywhere" binding:"required"`
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
		Name       string `json:"name" binding:"required"`
		StartMoney string `json:"startmoney" binding:"required"`
	}

	var req RconGameConfigStartMoneyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_startmoney "+req.StartMoney)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_startmoney " + req.StartMoney + " 响应: " + response)
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
		Name     string `json:"name" binding:"required"`
		MaxMoney string `json:"maxmoney" binding:"required"`
	}

	var req RconGameConfigMaxMoneyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_maxmoney "+req.MaxMoney)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_maxmoney " + req.MaxMoney + " 响应: " + response)
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
		Value string `json:"autoteambalance" binding:"required"`
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

// rconGameConfigLimitTeamsHandler 设置队伍人数限制
func rconGameConfigLimitTeamsHandler(c *gin.Context) {
	// 定义请求参数结构体
	type RconGameConfigLimitTeamsRequest struct {
		Name  string `json:"name" binding:"required"`
		Value string `json:"limitteams" binding:"required"`
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
		Name   string `json:"name" binding:"required"`
		C4Time string `json:"c4timer" binding:"required"`
	}

	var req RconGameConfigC4TimerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	response, err := ExecRconCommand(FullName(req.Name), "mp_c4timer "+req.C4Time)
	if err != nil {
		handleErrorResponse(c, "执行命令失败", err)
		return
	} else {
		util.Info("执行命令成功 命令: mp_c4timer " + req.C4Time + " 响应: " + response)
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
