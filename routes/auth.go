package routes

import (
	"github.com/galiherlangga/clinic-appointment/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, c *RouterContainer) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", c.AuthHandler.Login)
		auth.POST("/register", c.AuthHandler.Register)
		auth.POST("/refresh", c.AuthHandler.Refresh)
		auth.POST("/logout", c.AuthHandler.Logout)
		
		protected := auth.Group("")
		protected.Use(middleware.AuthMiddleware(c.BlacklistService))
		{
			protected.GET("/me", c.AuthHandler.Me)
		}
	}
}
