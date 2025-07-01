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
		// 域名管理路由组
		domainMgmt := api.Group("/domains")
		{
			// 主域名操作
			domainMgmt.GET("", dnsHandler.ListDomains) // 获取所有域名列表
			// TODO: 后续可以添加其他主域名相关操作，如：
			// - 添加域名
			// - 删除域名
			// - 修改域名分组
			// - 获取域名信息
		}

		// 解析记录管理路由组
		recordMgmt := domainMgmt.Group("/:domain/records")
		{
			// 基础操作
			recordMgmt.GET("", dnsHandler.ListDomainRecords)          // 获取域名的所有解析记录
			recordMgmt.GET("/search", dnsHandler.SearchDomainRecords) // 搜索解析记录（支持多条件）

			// 按标识符查询
			idQuery := recordMgmt.Group("/id")
			{
				idQuery.GET("/:record_id", dnsHandler.SearchDomainRecordsByRecordId) // 按记录ID查询
			}

			// 按记录属性查询
			attrQuery := recordMgmt.Group("")
			{
				attrQuery.GET("/rr/:rr", dnsHandler.SearchDomainRecordsByRR)             // 按主机记录查询
				attrQuery.GET("/type/:type", dnsHandler.SearchDomainRecordsByType)       // 按记录类型查询
				attrQuery.GET("/status/:status", dnsHandler.SearchDomainRecordsByStatus) // 按记录状态查询
			}

			// TODO: 后续可以添加解析记录的修改操作，如：
			// - 添加解析记录
			// - 修改解析记录
			// - 删除解析记录
			// - 设置解析记录状态
		}
	}

	return r
}
