package routes

import (
	"github.com/galiherlangga/clinic-appointment/handlers"
	"github.com/galiherlangga/clinic-appointment/middleware"
	"github.com/galiherlangga/clinic-appointment/services"
	"github.com/gin-gonic/gin"
)

type RouterContainer struct {
	AuthHandler        *handlers.AuthHandler
	CustomerHandler    *handlers.CustomerHandler
	ProviderHandler    *handlers.ProviderHandler
	AppointmentHandler *handlers.AppointmentHandler
	BlacklistService   services.BlacklistService
}

func SetupRoutes(r *gin.Engine, c *RouterContainer) {
	v1 := r.Group("/api/v1")
	v1.Use(middleware.SecurityMiddleware())
	{
		// Delegate to module-specific routing registrations
		RegisterAuthRoutes(v1, c)
		RegisterCustomerRoutes(v1, c)
		RegisterProviderRoutes(v1, c)
		RegisterAppointmentRoutes(v1, c)
		RegisterAdminRoutes(v1, c)
	}
}
