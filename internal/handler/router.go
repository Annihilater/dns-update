package handler

import (
	"dns-update/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter 初始化路由配置
func InitRouter(dnsHandler *DNSHandler) *gin.Engine {
	// 设置生产模式
	gin.SetMode(gin.ReleaseMode)

	// 创建 Gin 路由
	r := gin.New()
	r.Use(gin.Recovery())

	// 初始化Swagger文档
	docs.SwaggerInfo.BasePath = "/api"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API 路由组
	api := r.Group("/api")
	{
		// 域名相关路由
		domains := api.Group("/domains")
		{
			domains.GET("", dnsHandler.ListDomains)

			// 域名记录相关路由
			records := domains.Group("/:domain/records")
			{
				records.GET("", dnsHandler.ListDomainRecords)
				records.GET("/search", dnsHandler.SearchDomainRecords)
				records.GET("/id/:record_id", dnsHandler.SearchDomainRecordsByRecordId)
				records.GET("/rr/:rr", dnsHandler.SearchDomainRecordsByRR)
				records.GET("/type/:type", dnsHandler.SearchDomainRecordsByType)
				records.GET("/status/:status", dnsHandler.SearchDomainRecordsByStatus)
			}
		}
	}

	return r
}
