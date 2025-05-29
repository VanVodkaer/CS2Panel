package server

import (
	"github.com/VanVodkaer/CS2Panel/config"
	"github.com/gin-gonic/gin"
)

// infoMapUpdateHandler 处理获取地图列表的更新请求
func infoMapUpdateHandler(c *gin.Context) {
	// 定义请求参数结构体
	type MapListRequest struct {
		Class string `form:"class"` // "current" 或 "former"，不带参数时获取所有地图
	}

	var req MapListRequest
	if err := c.BindQuery(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}

	if req.Class == "current" {
		err := fetchCurrentMaps()
		if err != nil {
			handleErrorResponse(c, "获取当前地图列表失败", err)
			return
		}
	} else if req.Class == "former" {
		err := fetchFormerMaps()
		if err != nil {
			handleErrorResponse(c, "获取历史地图列表失败", err)
			return
		}
	} else {
		err := fetchCurrentMaps()
		if err != nil {
			handleErrorResponse(c, "获取当前地图列表失败", err)
			return
		}
		err = fetchFormerMaps()
		if err != nil {
			handleErrorResponse(c, "获取历史地图列表失败", err)
			return
		}
	}

	c.JSON(200, gin.H{
		"message": "地图列表更新成功",
	})
}

// infoMapListHandler 处理获取地图列表的请求
func infoMapListHandler(c *gin.Context) {
	// 定义请求参数结构体
	type MapListRequest struct {
		Class string `form:"class"`
	}

	var req MapListRequest
	if c.Request.ContentLength != 0 {
		if err := c.BindQuery(&req); err != nil {
			handleErrorResponse(c, "无效的请求参数", err)
			return
		}
	}

	if req.Class == "current" {
		maps, err := getCurrentMaps()
		if err != nil {
			handleErrorResponse(c, "获取当前地图列表失败", err)
			return
		}
		c.JSON(200, gin.H{
			"maps": maps,
		})
	} else if req.Class == "former" {
		maps, err := getFormerMaps()
		if err != nil {
			handleErrorResponse(c, "获取历史地图列表失败", err)
			return
		}
		c.JSON(200, gin.H{
			"maps": maps,
		})
	} else {
		var allMaps []MapInfo
		currentMaps, err := getCurrentMaps()
		if err != nil {
			handleErrorResponse(c, "获取当前地图列表失败", err)
			return
		}
		allMaps = append(allMaps, currentMaps...)
		formerMaps, err := getFormerMaps()
		if err != nil {
			handleErrorResponse(c, "获取历史地图列表失败", err)
			return
		}
		allMaps = append(allMaps, formerMaps...)
		c.JSON(200, gin.H{
			"maps": allMaps,
		})
	}
}

// infoNetworkAddrHandler 处理获取网络地址的请求
func infoNetworkAddrHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"addr": config.GlobalConfig.Game.Address,
	})
}

// infoNetworkGamePortHandler 处理获取网络端口的请求
func infoNetworkGamePortHandler(c *gin.Context) {
	// 定义请求参数结构体
	type NetworkPortRequest struct {
		Name string `form:"name" binding:"required"` // 容器名称
	}

	var req NetworkPortRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	// 获取网络端口信息
	port, err := GetEnvValue(FullName(req.Name), "CS2_PORT")
	if err != nil {
		handleErrorResponse(c, "获取网络端口失败", err)
		return
	}

	c.JSON(200, gin.H{
		"port": port,
	})
}

// infoNetworkTVPortHandler 处理获取网络端口的请求
func infoNetworkTVPortHandler(c *gin.Context) {
	// 定义请求参数结构体
	type NetworkPortRequest struct {
		Name string `form:"name" binding:"required"` // 容器名称
	}

	var req NetworkPortRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	// 获取网络端口信息
	port, err := GetEnvValue(FullName(req.Name), "TV_PORT")
	if err != nil {
		handleErrorResponse(c, "获取网络端口失败", err)
		return
	}

	c.JSON(200, gin.H{
		"port": port,
	})
}

// infoNetworkGamePasswdHandler 处理获取游戏密码的请求
func infoNetworkGamePasswdHandler(c *gin.Context) {
	// 定义请求参数结构体
	type NetworkPasswdRequest struct {
		Name string `form:"name" binding:"required"` // 容器名称
	}

	var req NetworkPasswdRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	// 获取游戏密码
	passwd, err := GetEnvValue(FullName(req.Name), "CS2_PW")
	if err != nil {
		handleErrorResponse(c, "获取游戏密码失败", err)
		return
	}

	c.JSON(200, gin.H{
		"passwd": passwd,
	})
}

// infoNetworkTVPasswdHandler 处理获取TV密码的请求
func infoNetworkTVPasswdHandler(c *gin.Context) {
	// 定义请求参数结构体
	type NetworkPasswdRequest struct {
		Name string `form:"name" binding:"required"` // 容器名称
	}

	var req NetworkPasswdRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		handleErrorResponse(c, "无效的请求参数", err)
		return
	}
	// 获取TV密码
	passwd, err := GetEnvValue(FullName(req.Name), "CS2_TV_PW")
	if err != nil {
		handleErrorResponse(c, "获取TV密码失败", err)
		return
	}

	c.JSON(200, gin.H{
		"passwd": passwd,
	})
}
