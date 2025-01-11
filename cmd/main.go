package main

import (
	"github.com/VanVodkaer/CS2Panel/server"
	"github.com/VanVodkaer/CS2Panel/util"
)

func main() {

	// 创建并初始化 Web 应用
	app, err := server.ServerNewApp()
	if err != nil {
		util.Error("初始化 Web 应用失败: %v", err)
	} else {
		util.Info("初始化 Web 应用成功")
	}
	// 启动 Web 服务器
	app.ServerStart()

}
