package middleware

import (
	"net/http"
	"strings"

	"go-biz-admin/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthMiddleware JWT验证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请求头中缺少Authorization字段"})
			c.Abort()
			return
		}

		// 验证格式是否为 "Bearer {token}"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header格式错误"})
			c.Abort()
			return
		}

		// 解析JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(services.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		// 从token中提取用户信息并放入上下文
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := uint(claims["user_id"].(float64))
			username, _ := claims["username"].(string)
			
			// 修复类型断言错误：检查role是否存在，否则默认为'user'
			role, exists := claims["role"]
			if !exists {
				role = "user"
			} else {
				if roleStr, ok := role.(string); ok {
					role = roleStr
				} else {
					role = "user"
				}
			}

			c.Set("user_id", userID)
			c.Set("username", username)
			c.Set("role", role)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminAuthMiddleware 管理员权限验证中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 首先确保已经通过JWT验证
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}
		
		c.Next()
	}
}