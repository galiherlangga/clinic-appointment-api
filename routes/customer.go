package routes

import (
	"github.com/galiherlangga/clinic-appointment/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterCustomerRoutes(rg *gin.RouterGroup, c *RouterContainer) {
	customers := rg.Group("/customers")
	customers.Use(middleware.AuthMiddleware(c.BlacklistService))
	{
		customers.GET("", c.CustomerHandler.FindAll)
		customers.GET("/:id", c.CustomerHandler.FindByID)
		customers.POST("", c.CustomerHandler.Create)
		customers.PUT("/:id", c.CustomerHandler.Update)
		customers.DELETE("/:id", c.CustomerHandler.Delete)
	}
}
