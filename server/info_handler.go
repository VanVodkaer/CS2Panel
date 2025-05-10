package server

import (
	"github.com/gin-gonic/gin"
)

type MapListRequest struct {
	Class string `json:"class"` // "current" 或 "former"，不带参数时获取所有地图
}

// rconMapListPostHandler 处理获取地图列表的更新请求
func infoMapUpdateHandler(c *gin.Context) {
	var req MapListRequest
	if c.Request.ContentLength != 0 {
		if err := c.BindJSON(&req); err != nil {
			handleErrorResponse(c, "无效的请求参数", err)
			return
		}
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

// rconMapListGetHandler 处理获取地图列表的请求
func infoMapListHandler(c *gin.Context) {
	var req MapListRequest
	if c.Request.ContentLength != 0 {
		if err := c.BindJSON(&req); err != nil {
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
