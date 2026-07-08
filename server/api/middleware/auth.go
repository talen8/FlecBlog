package middleware

import (
	"crypto/subtle"
	"strings"

	"flec_blog/config"
	"flec_blog/internal/service"
	"flec_blog/pkg/errcode"
	mcpauth "flec_blog/pkg/mcp/auth"
	"flec_blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// Auth JWT认证中间件
// 验证请求头中的 Bearer Token，并将用户信息注入上下文
// 同时支持从 URL 参数 token 中获取（用于 OAuth 绑定等浏览器跳转场景）
// 使用: router.Use(middleware.Auth(userService))
// 上下文注入: c.Get("user") 和 c.Get("user_id")
func Auth(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		// 优先从 Authorization header 获取
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		// 如果 header 中没有，尝试从 URL 参数获取（用于 OAuth 绑定跳转）
		if token == "" {
			token = c.Query("token")
		}

		if token == "" {
			response.Error(c, errcode.Unauthorized.WithDetails("未提供认证令牌"))
			c.Abort()
			return
		}

		user, err := userService.ValidateToken(token)
		if err != nil {
			response.Error(c, errcode.TokenError.WithDetails(err.Error()))
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Next()
	}
}

// OptionalAuth 可选JWT认证中间件
// 如果请求头中有有效的 Bearer Token，则将用户信息注入上下文
// 如果没有或无效，则继续执行，不中断请求（适用于支持游客和登录用户的接口）
// 使用: router.Use(middleware.OptionalAuth(userService))
// 上下文注入: c.Get("user") 和 c.Get("user_id") (可能为nil)
func OptionalAuth(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 没有 token，继续执行（作为游客）
			c.Next()
			return
		}

		// 检查 Bearer token 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			// token 格式无效，继续执行（作为游客）
			c.Next()
			return
		}

		// 验证 token
		user, err := userService.ValidateToken(parts[1])
		if err != nil {
			// token 无效，继续执行（作为游客）
			c.Next()
			return
		}

		// token 有效，设置用户信息
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Next()
	}
}

// MCPAuth MCP 专用 Bearer 鉴权中间件。
func MCPAuth(conf *config.Config) gin.HandlerFunc {
	return MCPAuthWithSecretProvider(func() string {
		if conf == nil {
			return ""
		}
		return conf.AI.MCPSecret
	})
}

// MCPAuthWithSecretProvider 使用动态 SecretProvider 验证旧版 MCP Secret。
func MCPAuthWithSecretProvider(secretProvider mcpauth.SecretProvider) gin.HandlerFunc {
	if secretProvider == nil {
		secretProvider = func() string { return "" }
	}
	return func(c *gin.Context) {
		token, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")
		secret := secretProvider()
		if ok && token != "" && secret != "" &&
			subtle.ConstantTimeCompare([]byte(token), []byte(secret)) == 1 {
			principal := mcpauth.StaticPrincipal()
			c.Request = c.Request.WithContext(mcpauth.ContextWithPrincipal(c.Request.Context(), principal))
			c.Set("mcp_principal", principal)
			c.Next()
			return
		}

		response.Error(c, errcode.Unauthorized.WithDetails("MCP 认证失败"))
		c.Abort()
	}
}

// MCPResourceAuth validates OAuth/hybrid bearer tokens for the MCP resource server.
func MCPResourceAuth(authenticator *mcpauth.Authenticator) gin.HandlerFunc {
	if authenticator == nil {
		panic("MCPResourceAuth requires a non-nil authenticator")
	}
	return func(c *gin.Context) {
		token, ok := mcpauth.ExtractBearerToken(c.GetHeader("Authorization"))
		if !ok {
			c.Header("WWW-Authenticate", authenticator.UnauthorizedChallenge())
			response.Error(c, errcode.Unauthorized.WithDetails("MCP 认证失败"))
			c.Abort()
			return
		}

		principal, err := authenticator.AuthenticateBearer(c.Request.Context(), token)
		if err != nil {
			c.Header("WWW-Authenticate", authenticator.UnauthorizedChallenge())
			response.Error(c, errcode.Unauthorized.WithDetails("MCP 认证失败"))
			c.Abort()
			return
		}

		c.Request = c.Request.WithContext(mcpauth.ContextWithPrincipal(c.Request.Context(), principal))
		c.Set("mcp_principal", principal)
		c.Next()
	}
}
