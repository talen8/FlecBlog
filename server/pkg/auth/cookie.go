package auth

import (
	"net/http"

	"flec_blog/pkg/utils"

	"github.com/gin-gonic/gin"
)

// SetRefreshTokenCookie 设置 refresh token 到 HttpOnly Cookie
func SetRefreshTokenCookie(ctx *gin.Context, refreshToken string) {
	ctx.SetSameSite(getSameSiteMode(ctx))
	ctx.SetCookie(
		"refresh_token",
		refreshToken,
		utils.RefreshTokenExpireHours*3600,
		"/",
		"",
		isHTTPS(ctx),
		true,
	)
}

// ClearRefreshTokenCookie 清除 refresh token Cookie
func ClearRefreshTokenCookie(ctx *gin.Context) {
	ctx.SetSameSite(getSameSiteMode(ctx))
	ctx.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		isHTTPS(ctx),
		true,
	)
}

// GetRefreshTokenFromCookie 从 Cookie 获取 refresh token
func GetRefreshTokenFromCookie(ctx *gin.Context) string {
	token, _ := ctx.Cookie("refresh_token")
	return token
}

// getSameSiteMode HTTPS 用 None（兼容跨域），HTTP 用 Lax（开发环境）
func getSameSiteMode(ctx *gin.Context) http.SameSite {
	if isHTTPS(ctx) {
		return http.SameSiteNoneMode
	}
	return http.SameSiteLaxMode
}

// isHTTPS 判断是否为 HTTPS
func isHTTPS(ctx *gin.Context) bool {
	if ctx.Request.TLS != nil {
		return true
	}
	return ctx.GetHeader("X-Forwarded-Proto") == "https"
}
