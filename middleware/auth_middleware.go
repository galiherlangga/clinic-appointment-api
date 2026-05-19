package middleware

import (
	"net/http"
	"strings"

	"github.com/galiherlangga/clinic-appointment/configs"
	"github.com/galiherlangga/clinic-appointment/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(blacklistService services.BlacklistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Check Blacklist
		isBlacklisted, _ := blacklistService.IsTokenBlacklisted(tokenString)
		if isBlacklisted {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is blacklisted (logged out)"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(configs.AppConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		userID := uint(claims["user_id"].(float64))
		role := claims["role"].(string)

		c.Set("user_id", userID)
		c.Set("role", role)

		c.Next()
	}
}
