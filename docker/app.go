package docker

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/VanVodkaer/CS2Panel/config"
	"github.com/VanVodkaer/CS2Panel/util"
	"github.com/docker/docker/client"
)

// 全局变量 Cli 为 Docker 客户端对象
var Cli *client.Client

// init 初始化 Docker 客户端对象
func init() {
	// 初始化 Docker 客户端对象
	var err error
	Cli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		util.Error("初始化 Docker 客户端失败", err)
		return
	} else {
		util.Info("初始化 Docker 客户端成功")
	}

	// 测试 Docker Daemon 连接（包含重试机制）
	if err := TestDockerConnection(); err != nil {
		util.Error("Docker Daemon 连接失败", err)
		os.Exit(1)
	}
}

// TestDockerConnection 封装 Docker Daemon 连接测试的逻辑
func TestDockerConnection() error {
	// 设置最大重试次数和间隔
	maxRetries := config.GlobalConfig.Docker.MaxRetries
	retryDelay := time.Duration(config.GlobalConfig.Docker.RetryDelay) * time.Second

	var err error
	for i := 0; i < maxRetries; i++ {
		_, err = Cli.Ping(context.Background())
		if err == nil {
			// 如果连接成功，返回 nil
			util.Debug("测试 Docker Daemon 连接成功")
			return nil
		}

		// 如果连接失败，打印日志并等待一段时间后重试
		util.Error(fmt.Sprintf("测试 Docker Daemon 连接失败, 重试 %d/%d", i+1, maxRetries), err)
		time.Sleep(retryDelay)
	}

	// 如果所有重试都失败，返回最终的错误
	return fmt.Errorf("测试 Docker Daemon 连接失败, 请检查 Docker 服务: %v", err)

}
