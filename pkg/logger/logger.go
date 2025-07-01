package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// InitLogger 初始化日志
func InitLogger() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	var err error
	Log, err = config.Build()
	if err != nil {
		panic("初始化日志失败: " + err.Error())
	}
}

// GetLogger 获取日志实例
func GetLogger() *zap.Logger {
	if Log == nil {
		InitLogger()
	}
	return Log
}
