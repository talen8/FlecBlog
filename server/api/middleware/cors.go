package middleware

import (
	"flec_blog/config"

	"github.com/gin-gonic/gin"
)

// CORS 跨域中间件
func CORS(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 检查白名单
		allowed := false
		for _, allowOrigin := range cfg.Server.AllowOrigins {
			if allowOrigin == "*" || allowOrigin == origin {
				allowed = true
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}
		if !allowed && origin != "" {
			c.AbortWithStatus(403)
			return
		}

		// 允许携带凭据（cookies、认证信息等）
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "X-Scene")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
