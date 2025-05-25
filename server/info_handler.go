package server

import (
	"fmt"

	"github.com/VanVodkaer/CS2Panel/config"
	"github.com/gin-gonic/gin"
)

type MapListRequest struct {
	Class string `form:"class"` // "current" 或 "former"，不带参数时获取所有地图
}

// infoMapUpdateHandler 处理获取地图列表的更新请求
func infoMapUpdateHandler(c *gin.Context) {
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

// networkAddrHandler 处理获取网络地址的请求
func networkAddrHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"addr": config.GlobalConfig.Game.Address,
	})
}

type NetworkPortRequest struct {
	Name string `form:"name" binding:"required"` // 容器名称
}

// networkPortHandler 处理获取网络端口的请求
func networkGamePortHandler(c *gin.Context) {
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

// networkTVPortHandler 处理获取网络端口的请求
func networkTVPortHandler(c *gin.Context) {
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

// networkGamePortsHandler 处理获取网络端口的请求
func networkGamePortsHandler(c *gin.Context) {
	// 使用 QueryArray 获取 name 参数列表
	names := c.QueryArray("name")
	if len(names) == 0 {
		handleErrorResponse(c, "无效的请求参数", fmt.Errorf("缺少 name 参数"))
		return
	}

	// 获取网络端口信息
	var allPorts []string
	for _, name := range names {
		ports, err := GetEnvValue(FullName(name), "CS2_PORT")
		if err != nil {
			handleErrorResponse(c, "获取网络端口失败", err)
			return
		}
		allPorts = append(allPorts, ports)
	}

	c.JSON(200, gin.H{
		"ports": allPorts,
	})
}

// networkTVPortsHandler 处理获取网络端口的请求
func networkTVPortsHandler(c *gin.Context) {
	names := c.QueryArray("name")
	if len(names) == 0 {
		handleErrorResponse(c, "无效的请求参数", fmt.Errorf("缺少 name 参数"))
		return
	}
	// 获取TV端口信息
	var allPorts []string
	for _, name := range names {
		ports, err := GetEnvValue(FullName(name), "TV_PORT")
		if err != nil {
			handleErrorResponse(c, "获取TV端口失败", err)
			return
		}
		allPorts = append(allPorts, ports)
	}

	c.JSON(200, gin.H{
		"ports": allPorts,
	})
}

type NetworkPasswdRequest struct {
	Name string `form:"name" binding:"required"` // 容器名称
}

// networkGamePasswdHandler 处理获取游戏密码的请求
func networkGamePasswdHandler(c *gin.Context) {
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

// networkTVPasswdHandler 处理获取TV密码的请求
func networkTVPasswdHandler(c *gin.Context) {
	var req NetworkPasswdRequest
	if c.Request.ContentLength != 0 {
		if err := c.ShouldBindQuery(&req); err != nil {
			handleErrorResponse(c, "无效的请求参数", err)
			return
		}
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

// networkGamePasswdsHandler 处理获取游戏密码的请求
func networkGamePasswdsHandler(c *gin.Context) {

	names := c.QueryArray("name")
	if len(names) == 0 {
		handleErrorResponse(c, "无效的请求参数", fmt.Errorf("缺少 name 参数"))
		return
	}
	// 获取游戏密码
	var allPasswds []string
	for _, name := range names {
		passwd, err := GetEnvValue(FullName(name), "CS2_PW")
		if err != nil {
			handleErrorResponse(c, "获取游戏密码失败", err)
			return
		}
		allPasswds = append(allPasswds, passwd)
	}

	c.JSON(200, gin.H{
		"passwds": allPasswds,
	})
}

// networkTVPasswdsHandler 处理获取TV密码的请求
func networkTVPasswdsHandler(c *gin.Context) {

	names := c.QueryArray("name")
	if len(names) == 0 {
		handleErrorResponse(c, "无效的请求参数", fmt.Errorf("缺少 name 参数"))
		return
	}
	// 获取TV密码
	var allPasswds []string
	for _, name := range names {
		passwd, err := GetEnvValue(FullName(name), "CS2_TV_PW")
		if err != nil {
			handleErrorResponse(c, "获取TV密码失败", err)
			return
		}
		allPasswds = append(allPasswds, passwd)
	}
	c.JSON(200, gin.H{
		"passwds": allPasswds,
	})
}
