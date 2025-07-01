package main

import (
	"fmt"
	"os"

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
	defer logger.Log.Sync()
	log := logger.GetLogger()

	// 获取环境变量
	accessKeyId := os.Getenv("ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ACCESS_KEY_SECRET")

	// 验证必要的环境变量
	if accessKeyId == "" || accessKeySecret == "" {
		log.Fatal("缺少必要的环境变量",
			zap.String("ACCESS_KEY_ID", accessKeyId),
			zap.String("ACCESS_KEY_SECRET", "***"),
		)
	}

	// 初始化 DNS 服务
	dnsService, err := service.NewDNSService(
		tea.String(accessKeyId),
		tea.String(accessKeySecret),
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

	// 获取服务端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 打印服务信息
	log.Info("DNS Update Service is running",
		zap.String("swagger_url", fmt.Sprintf("http://localhost:%s/swagger/index.html", port)),
		zap.String("server_url", fmt.Sprintf("http://localhost:%s", port)),
	)

	// 启动服务器
	if err := r.Run(":" + port); err != nil {
		log.Fatal("启动服务失败", zap.Error(err))
	}
}
