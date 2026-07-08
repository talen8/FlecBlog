package middleware

import (
	"net/http"

	"flec_blog/config"

	"github.com/gin-gonic/gin"
)

// CORS 跨域中间件
func CORS(cfg *config.Config) gin.HandlerFunc {
	return CORSWithNavigationPaths(cfg)
}

// CORSWithNavigationPaths keeps the normal CORS allowlist while permitting
// browser navigation to selected exact paths. GET navigations are public;
// POST navigations must be browser-reported same-origin document navigations.
// A navigation exception does not grant CORS read access: no
// Access-Control-Allow-Origin header is emitted for an unallowlisted origin.
func CORSWithNavigationPaths(cfg *config.Config, navigationPaths ...string) gin.HandlerFunc {
	navigationPathSet := make(map[string]struct{}, len(navigationPaths))
	for _, path := range navigationPaths {
		navigationPathSet[path] = struct{}{}
	}

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

		navigationException := false
		if !allowed && origin != "" {
			_, exactNavigationPath := navigationPathSet[c.Request.URL.Path]
			navigationException = exactNavigationPath && isAllowedNavigationRequest(c.Request)
			if !navigationException {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		if !navigationException {
			// 允许携带凭据（cookies、认证信息等）
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Writer.Header().Set("Access-Control-Expose-Headers", "X-Scene")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func isAllowedNavigationRequest(r *http.Request) bool {
	switch r.Method {
	case http.MethodGet:
		return true
	case http.MethodPost:
		return r.Header.Get("Sec-Fetch-Mode") == "navigate" &&
			r.Header.Get("Sec-Fetch-Dest") == "document" &&
			r.Header.Get("Sec-Fetch-Site") == "same-origin"
	default:
		return false
	}
}
