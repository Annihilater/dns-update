package main

import (
	"log"
	"os"

	"dns-update/internal/handler"
	"dns-update/internal/service"

	env "github.com/alibabacloud-go/darabonba-env/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
)

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

	// 创建 Gin 路由
	r := gin.Default()

	// 注册路由
	r.GET("/api/domains", dnsHandler.ListDomains)
	r.GET("/api/domains/:domain/records", dnsHandler.ListDomainRecords)

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
