package main

import (
	"fmt"
	"log"
	"os"

	"dns-update/docs"
	"dns-update/internal/handler"
	"dns-update/internal/service"

	env "github.com/alibabacloud-go/darabonba-env/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title        DNS Update API
// @version      1.0
// @description  阿里云DNS管理服务API
// @BasePath     /api

func main() {
	// 初始化 DNS 服务
	dnsService, err := service.NewDNSService(
		env.GetEnv(tea.String("ACCESS_KEY_ID")),
		env.GetEnv(tea.String("ACCESS_KEY_SECRET")),
	)
	if err != nil {
		log.Fatalf("Failed to initialize DNS service: %v", err)
	}

	// 初始化处理器
	dnsHandler := handler.NewDNSHandler(dnsService)

	// 设置生产模式
	gin.SetMode(gin.ReleaseMode)

	// 创建 Gin 路由
	r := gin.New()
	r.Use(gin.Recovery())

	// 初始化Swagger文档
	docs.SwaggerInfo.BasePath = "/api"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 注册路由
	r.GET("/api/domains", dnsHandler.ListDomains)
	r.GET("/api/domains/:domain/records", dnsHandler.ListDomainRecords)

	// 获取服务端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 打印服务信息
	fmt.Printf("\n📡 DNS Update Service is running:\n")
	fmt.Printf("🌐 API Documentation: http://localhost:%s/swagger/index.html\n", port)
	fmt.Printf("🚀 HTTP Server: http://localhost:%s\n\n", port)

	// 启动服务器
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
