package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/config"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
	"github.com/wxlbd/nutri-baby-server/pkg/response"
)

// Auth JWT认证中间件
func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, errors.ErrUnauthorized)
			c.Abort()
			return
		}

		// 验证Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, errors.ErrUnauthorized)
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析Token
		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			response.Error(c, errors.ErrInvalidToken)
			c.Abort()
			return
		}

		// 获取Claims
		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok {
			response.Error(c, errors.ErrInvalidToken)
			c.Abort()
			return
		}

		// 设置用户信息到context
		c.Set("openid", claims.Subject)

		c.Next()
	}
}
