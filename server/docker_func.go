package server

import (
	"github.com/VanVodkaer/CS2Panel/config"
)

// FullName 生成完整的名称 prefix-name
func FullName(name string) string {
	return config.GlobalConfig.Docker.Prefix + "-" + name
}
