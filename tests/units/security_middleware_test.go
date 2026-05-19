package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/galiherlangga/clinic-appointment/configs"
	"github.com/galiherlangga/clinic-appointment/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSecurityMiddleware(t *testing.T) {
	// Setup config for testing
	configs.AppConfig = &configs.Config{
		APIKey:     "test-api-key",
		DevKey:     "test-dev-key",
		AllowedIPs: []string{"1.2.3.4"},
	}

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		apiKey         string
		devKey         string
		clientIP       string
		expectedStatus int
	}{
		{
			name:           "Missing API Key",
			apiKey:         "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid API Key",
			apiKey:         "wrong-key",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Valid API Key, Non-Whitelisted IP",
			apiKey:         "test-api-key",
			clientIP:       "5.6.7.8",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Valid API Key, Whitelisted IP",
			apiKey:         "test-api-key",
			clientIP:       "1.2.3.4",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid API Key, Dev Key Bypass (Any IP)",
			apiKey:         "test-api-key",
			devKey:         "test-dev-key",
			clientIP:       "9.9.9.9",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(func(c *gin.Context) {
				if tt.clientIP != "" {
					// Mock client IP
					c.Request.RemoteAddr = tt.clientIP + ":1234"
				}
				c.Next()
			})
			r.Use(middleware.SecurityMiddleware())
			r.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			if tt.apiKey != "" {
				req.Header.Set("X-API-KEY", tt.apiKey)
			}
			if tt.devKey != "" {
				req.Header.Set("X-DEV-KEY", tt.devKey)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
