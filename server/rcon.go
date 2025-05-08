package server

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/VanVodkaer/CS2Panel/config"
	"github.com/VanVodkaer/CS2Panel/docker"
	"github.com/VanVodkaer/CS2Panel/util"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	rcon "github.com/forewing/csgo-rcon"
)

// ExecRconCommand 执行Rcon命令
func ExecRconCommand(name string, command string) (string, error) {
	port := GetRconPort(name)
	if port == 0 {
		return "", fmt.Errorf("获取Rcon端口失败，请检查容器是否存在或Rcon端口是否正确配置")
	}
	// 创建Rcon客户端
	client := rcon.New(fmt.Sprintf("localhost:%d", port), config.GlobalConfig.Game.RCON_PASSWORD, 1*time.Second)

	response, err := client.Execute(command)
	if err != nil {
		return "", fmt.Errorf("执行Rcon命令失败: %v", err)
	}

	return response, err
}

// GetRconPort 获取Rcon端口
func GetRconPort(name string) int {
	// 获取所有容器
	containers, err := docker.Cli.ContainerList(context.Background(), container.ListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("name", config.GlobalConfig.Docker.Prefix)),
	})
	if err != nil {
		util.Error("获取 Docker 容器列表失败:", err)
		return 0
	}

	for _, c := range containers {
		// 检查容器ID是否匹配
		if strings.Contains(c.Names[0], name) {
			// 遍历容器的端口映射
			tcpPorts := make(map[uint16]bool)
			udpPorts := make(map[uint16]bool)

			fmt.Printf("%+v\n", c.Ports)

			for _, port := range c.Ports {
				switch port.Type {
				case "tcp":
					tcpPorts[port.PublicPort] = true
				case "udp":
					udpPorts[port.PublicPort] = true
				}
			}

			// 找到同时存在于 TCP 和 UDP 的端口
			for port := range tcpPorts {
				if udpPorts[port] {
					return int(port)
				}
			}
		}
	}

	// 未找到则返回0
	return 0
}
