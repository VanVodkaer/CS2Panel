package server

import (
	"fmt"
	"strconv"
	"time"

	rcon "github.com/forewing/csgo-rcon"
)

// ExecRconCommand 执行Rcon命令
func ExecRconCommand(name string, command string) (string, error) {
	port, err := GetRconPort(name)
	if err != nil {
		return "", fmt.Errorf("获取Rcon端口失败，请检查容器是否存在或Rcon端口是否正确配置: %v", err)
	}
	// 创建Rcon客户端
	passwd, err := GetEnvValue(name, "CS2_RCONPW")
	if err != nil {
		return "", fmt.Errorf("获取Rcon密码失败: %v", err)
	}
	client := rcon.New(fmt.Sprintf("localhost:%d", port), passwd, 1*time.Second)

	// 执行Rcon命令
	response, err := client.Execute(command)
	if err != nil {
		return "", fmt.Errorf("执行Rcon命令失败: %v", err)
	}

	return response, nil
}

func ExecRconCommands(name string, commands []string) ([]string, error) {
	var responses []string
	for _, cmd := range commands {
		res, err := ExecRconCommand(name, cmd)
		if err != nil {
			return nil, err
		}
		responses = append(responses, res)
	}
	return responses, nil
}

// GetRconPort 获取Rcon端口
func GetRconPort(name string) (int, error) {
	if port, err := GetEnvValue(name, "RCON_PORT"); err == nil {
		return strconv.Atoi(port)
	}
	return 0, fmt.Errorf("未找到环境变量 %q", "RCON_PORT")
}
