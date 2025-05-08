package util

import "github.com/VanVodkaer/CS2Panel/config"

func DefaultIfEmpty[T comparable](value, defaultValue T) T {
	var zeroValue T
	if value == zeroValue {
		return defaultValue
	}
	return value
}

func FullName(name string) string {
	return config.GlobalConfig.Docker.Prefix + "-" + name
}
