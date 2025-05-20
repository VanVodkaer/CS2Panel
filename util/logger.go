package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/VanVodkaer/CS2Panel/config"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

// Logger 全局日志记录器
var Logger *logrus.Logger

// 初始化日志记录器
func init() {

	// 获取日志配置
	logDir := config.GlobalConfig.Util.LogDir
	logFileName := config.GlobalConfig.Util.LogFileName
	logMaxSize := config.GlobalConfig.Util.LogMaxSize
	logMaxBackups := config.GlobalConfig.Util.LogMaxBackups
	logMaxAge := config.GlobalConfig.Util.LogMaxAge
	logCompress := config.GlobalConfig.Util.LogCompress
	logLevel := config.GlobalConfig.Env.LogLevel
	// 日志文件路径
	logFilePath := filepath.Join(logDir, logFileName)

	// 创建日志记录器实例
	Logger = logrus.New()

	// 检查并创建日志目录
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("无法创建日志目录 '%s': %v\n", logDir, err)
		return
	}

	// 设置日志文件轮换
	Logger.Out = &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    logMaxSize,
		MaxBackups: logMaxBackups,
		MaxAge:     logMaxAge,
		Compress:   logCompress,
	}

	// 设置日志格式为文本格式
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// 设置日志级别
	var level logrus.Level
	switch logLevel {
	case "error":
		level = logrus.ErrorLevel
	case "warn":
		level = logrus.WarnLevel
	case "info":
		level = logrus.InfoLevel
	case "debug":
		level = logrus.DebugLevel
	default:
		logrus.Warnf("错误的日志级别: %s, 设置默认级别Info", level)
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

}

// Info 记录信息级别日志
func Info(msg string) {
	if Logger != nil {
		Logger.Info(msg)
		if config.GlobalConfig.Env.Mode == "debug" {
			log.Println(msg)
		}
	} else {
		logrus.Info(msg)
	}
}

// Error 记录错误级别日志
func Error(msg string, err error) {
	if Logger != nil {
		Logger.WithError(err).Error(msg)
		if config.GlobalConfig.Env.Mode == "debug" {
			log.Println(msg)
		}
	} else {
		logrus.WithError(err).Error(msg)
	}
}

// Warn 记录警告级别日志
func Warn(msg string) {
	if Logger != nil {
		Logger.Warn(msg)
		if config.GlobalConfig.Env.Mode == "debug" {
			log.Println(msg)
		}
	} else {
		logrus.Warn(msg)
	}
}

// Debug 记录调试级别日志
func Debug(msg string) {
	if Logger != nil {
		Logger.Debug(msg)
		if config.GlobalConfig.Env.Mode == "debug" {
			log.Println(msg)
		}
	} else {
		logrus.Debug(msg)
	}
}
