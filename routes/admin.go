package routes

import (
	"github.com/galiherlangga/clinic-appointment/middleware"
	"github.com/galiherlangga/clinic-appointment/models"
	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(rg *gin.RouterGroup, c *RouterContainer) {
	admin := rg.Group("/admin")
	admin.Use(middleware.AuthMiddleware(c.BlacklistService), middleware.RBACMiddleware(models.RoleAdmin, models.RoleSuperAdmin))
	{
		admin.GET("/dashboard", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"message": "welcome to admin dashboard"})
		})
	}
}
