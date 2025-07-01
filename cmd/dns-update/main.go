package main

import (
	"fmt"

	"dns-update/internal/config"
	"dns-update/internal/handler"
	"dns-update/internal/middleware"
	"dns-update/internal/service"
	"dns-update/pkg/logger"

	"github.com/alibabacloud-go/tea/tea"
	"go.uber.org/zap"
)

// @title        DNS Update API
// @version      1.0
// @description  阿里云DNS管理服务API
// @BasePath     /api

func main() {
	// 初始化日志
	logger.InitLogger()
	defer func(Log *zap.Logger) {
		err := Log.Sync()
		if err != nil {
			logger.GetLogger().Error("日志同步失败", zap.Error(err))
		}
	}(logger.Log)
	log := logger.GetLogger()

	// 加载配置
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatal("加载配置失败", zap.Error(err))
	}

	// 初始化 DNS 服务
	dnsService, err := service.NewDNSService(
		tea.String(cfg.Aliyun.AccessKeyId),
		tea.String(cfg.Aliyun.AccessKeySecret),
	)
	if err != nil {
		log.Fatal("初始化DNS服务失败", zap.Error(err))
	}

	// 初始化处理器
	dnsHandler := handler.NewDNSHandler(dnsService)

	// 初始化路由
	r := handler.InitRouter(dnsHandler)

	// 添加中间件
	r.Use(middleware.RequestTimer())

	// 打印服务信息
	log.Info("DNS Update Service is running",
		zap.String("swagger_url", fmt.Sprintf("http://localhost:%s/swagger/index.html", cfg.Server.Port)),
		zap.String("server_url", fmt.Sprintf("http://localhost:%s", cfg.Server.Port)),
	)

	// 启动服务器
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("启动服务失败", zap.Error(err))
	}
}
