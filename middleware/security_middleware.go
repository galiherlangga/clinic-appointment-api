package middleware

import (
	"net/http"

	"github.com/galiherlangga/clinic-appointment/configs"
	"github.com/gin-gonic/gin"
)

func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")
		devKey := c.GetHeader("X-DEV-KEY")
		clientIP := c.ClientIP()

		// 1. API Key Check (Mandatory for all)
		if apiKey == "" || apiKey != configs.AppConfig.APIKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or missing API Key"})
			c.Abort()
			return
		}

		// 2. Dev Key Bypass (If valid Dev Key is provided, skip IP check)
		if devKey != "" && devKey == configs.AppConfig.DevKey {
			c.Next()
			return
		}

		// 3. IP Whitelisting Check
		isAllowed := false
		for _, ip := range configs.AppConfig.AllowedIPs {
			if clientIP == ip {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: IP not whitelisted"})
			c.Abort()
			return
		}

		c.Next()
	}
}
