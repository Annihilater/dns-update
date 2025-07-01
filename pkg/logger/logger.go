package logger

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// LogConfig 日志配置
type LogConfig struct {
	LogPath    string // 日志文件路径
	MaxSize    int    // 每个日志文件的最大大小（MB）
	MaxBackups int    // 保留的旧日志文件的最大数量
	MaxAge     int    // 保留的旧日志文件的最大天数
	Compress   bool   // 是否压缩旧日志文件
}

// DefaultLogConfig 默认日志配置
var DefaultLogConfig = LogConfig{
	LogPath:    "logs/app.log",
	MaxSize:    100,
	MaxBackups: 30,
	MaxAge:     7,
	Compress:   true,
}

// InitLogger 初始化日志
func InitLogger() {
	initLoggerWithConfig(DefaultLogConfig)
}

// initLoggerWithConfig 使用指定配置初始化日志
func initLoggerWithConfig(config LogConfig) {
	// 确保日志目录存在
	logDir := filepath.Dir(config.LogPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic("创建日志目录失败: " + err.Error())
	}

	// 创建编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建文件输出
	fileWriter, err := os.OpenFile(config.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("打开日志文件失败: " + err.Error())
	}

	// 创建Core
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(fileWriter),
		zapcore.InfoLevel,
	)

	// 创建控制台输出
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)

	// 合并多个Core
	core := zapcore.NewTee(fileCore, consoleCore)

	// 创建Logger
	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

// GetLogger 获取日志实例
func GetLogger() *zap.Logger {
	if Log == nil {
		InitLogger()
	}
	return Log
}

// RotateLogFile 轮转日志文件
func RotateLogFile() error {
	// 获取当前时间
	now := time.Now()
	timestamp := now.Format("20060102150405")

	// 构建新的文件名
	logPath := DefaultLogConfig.LogPath
	dir := filepath.Dir(logPath)
	base := filepath.Base(logPath)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]
	rotatedPath := filepath.Join(dir, name+"."+timestamp+ext)

	// 关闭当前日志文件
	if Log != nil {
		Log.Sync()
	}

	// 重命名当前日志文件
	if err := os.Rename(logPath, rotatedPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	// 重新初始化日志
	InitLogger()

	return nil
}
