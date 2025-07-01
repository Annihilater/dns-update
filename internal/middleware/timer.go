package middleware

import (
	"dns-update/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestTimer 请求耗时统计中间件
func RequestTimer() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取开始时间
		start := time.Now()

		// 处理请求
		c.Next()

		// 计算耗时
		duration := time.Since(start)

		// 将耗时写入响应头（毫秒为单位）
		c.Header("X-Response-Time", duration.String())

		// 写入日志
		logger.GetLogger().Info("请求处理完成",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
		)
	}
}
