package docker

import (
	"context"

	"github.com/VanVodkaer/CS2Panel/config"
	"github.com/VanVodkaer/CS2Panel/util"
	"github.com/docker/docker/client"
)

// 全局变量 Cli 为 Docker 客户端对象
var Cli *client.Client

func init() {
	// 初始化 Docker 客户端对象
	Cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		util.Error("初始化 Docker 客户端失败", err)
	} else {
		util.Info("初始化 Docker 客户端成功")
	}

	if config.GlobalConfig.Env.Mode == "debug" {
		// 测试 Docker Daemon 连接
		_, err = Cli.Ping(context.Background())
		if err != nil {
			util.Error("测试 Docker Daemon 连接失败", err)
		} else {
			util.Debug("测试 Docker Daemon 连接成功")
		}
	}

	return
}
