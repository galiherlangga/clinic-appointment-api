package routes

import (
	"github.com/galiherlangga/clinic-appointment/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterProviderRoutes(rg *gin.RouterGroup, c *RouterContainer) {
	providers := rg.Group("/providers")
	providers.Use(middleware.AuthMiddleware(c.BlacklistService))
	{
		providers.GET("", c.ProviderHandler.FindAll)
		providers.GET("/:id", c.ProviderHandler.FindByID)
		providers.POST("", c.ProviderHandler.Create)
		providers.PUT("/:id", c.ProviderHandler.Update)
		providers.DELETE("/:id", c.ProviderHandler.Delete)
	}
}
