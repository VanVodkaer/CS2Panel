package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/VanVodkaer/CS2Panel/config"
	"github.com/VanVodkaer/CS2Panel/docker"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

// FullName 生成完整的名称 prefix-name
func FullName(name string) string {
	return config.GlobalConfig.Docker.Prefix + "-" + name
}

func GetEnvValue(name string, key string) (string, error) {
	// 获取所有容器
	containers, err := docker.Cli.ContainerList(context.Background(), container.ListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("name", config.GlobalConfig.Docker.Prefix)),
	})
	if err != nil {
		return "", err
	}

	for _, c := range containers {
		// 检查容器ID是否匹配
		if strings.TrimPrefix(c.Names[0], "/") == name {
			// 获取容器详细信息
			containerInfo, err := docker.Cli.ContainerInspect(context.Background(), c.ID)
			if err != nil {
				return "", err
			}

			// 打印容器的环境变量
			for _, env := range containerInfo.Config.Env {
				if strings.HasPrefix(env, key) {
					// 解析环境变量 通过=分割
					parts := strings.SplitN(env, "=", 2)
					if len(parts) == 2 {
						return parts[1], nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("未找到环境变量 %q", key)
}
