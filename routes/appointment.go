package routes

import (
	"github.com/galiherlangga/clinic-appointment/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAppointmentRoutes(rg *gin.RouterGroup, c *RouterContainer) {
	appointments := rg.Group("/appointments")
	appointments.Use(middleware.AuthMiddleware(c.BlacklistService))
	{
		appointments.GET("", c.AppointmentHandler.FindAll)
		appointments.GET("/:id", c.AppointmentHandler.FindByID)
		appointments.POST("", c.AppointmentHandler.Create)
		appointments.PUT("/:id", c.AppointmentHandler.Update)
		appointments.DELETE("/:id", c.AppointmentHandler.Delete)
	}
}
