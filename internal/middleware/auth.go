package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"PLMS/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Authorization
		authHeader := c.GetHeader("Authorization")

		// 调试日志：打印请求信息
		fmt.Printf("[Auth] Request: %s %s, Authorization: %s\n", c.Request.Method, c.Request.URL.Path, authHeader)

		if authHeader == "" {
			fmt.Printf("[Auth] 401: 请求头中没有 Authorization\n")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未登录或token无效",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "token格式错误",
				"data":    nil,
			})
			c.Abort()
			return
		}

		token := parts[1]
		fmt.Printf("[Auth] Token: %s...\n", token[:min(30, len(token))])

		// 验证token
		claims, err := services.ParseToken(token)
		if err != nil {
			fmt.Printf("[Auth] 401: Token解析失败: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "token无效或已过期: " + err.Error(),
				"data":    nil,
			})
			c.Abort()
			return
		}

		fmt.Printf("[Auth] Token验证成功, userID: %d, username: %s\n", claims.UserID, claims.Username)

		// 将用户信息存入上下文
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// DefaultPasswordCheckMiddleware 检查是否使用默认密码
// 如果使用默认密码，只允许访问修改密码和登录相关接口
func DefaultPasswordCheckMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.Next()
			return
		}

		// 检查当前请求是否是修改密码或登出接口
		path := c.Request.URL.Path
		if path == "/api/v1/auth/change-password" ||
			path == "/api/v1/auth/logout" ||
			path == "/api/v1/auth/current-user" ||
			path == "/api/v1/auth/check-default-password" {
			c.Next()
			return
		}

		// 获取用户服务检查是否使用默认密码
		authService := services.NewAuthService(db)
		isDefault, err := authService.IsDefaultPassword(userID.(int64))
		if err != nil {
			c.Next()
			return
		}

		if isDefault {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "请先修改默认密码",
				"data": gin.H{
					"requireChangePassword": true,
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// min 返回较小的整数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// CORSMiddleware 跨域中间件（增强版）
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}

		// 设置跨域头
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, token, x-requested-with, X-Token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Vary", "Origin")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
