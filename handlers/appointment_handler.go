package handlers

import (
	"net/http"
	"strconv"

	"github.com/galiherlangga/clinic-appointment/models"
	"github.com/galiherlangga/clinic-appointment/services"
	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	appointmentService services.AppointmentService
}

func NewAppointmentHandler(appointmentService services.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{appointmentService: appointmentService}
}

func (h *AppointmentHandler) FindAll(c *gin.Context) {
	var filter models.AppointmentFilter
	
	if custIDStr := c.Query("customer_id"); custIDStr != "" {
		if custID, err := strconv.ParseUint(custIDStr, 10, 32); err == nil {
			filter.CustomerID = uint(custID)
		}
	}
	if provIDStr := c.Query("provider_id"); provIDStr != "" {
		if provID, err := strconv.ParseUint(provIDStr, 10, 32); err == nil {
			filter.ProviderID = uint(provID)
		}
	}
	if svcIDStr := c.Query("service_id"); svcIDStr != "" {
		if svcID, err := strconv.ParseUint(svcIDStr, 10, 32); err == nil {
			filter.ServiceID = uint(svcID)
		}
	}
	filter.Status = c.Query("status")

	appointments, err := h.appointmentService.GetAll(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointments)
}

func (h *AppointmentHandler) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	appointment, err := h.appointmentService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "appointment not found"})
		return
	}

	c.JSON(http.StatusOK, appointment)
}

func (h *AppointmentHandler) Create(c *gin.Context) {
	var appointment models.Appointment
	if err := c.ShouldBindJSON(&appointment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.appointmentService.Create(&appointment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, appointment)
}

func (h *AppointmentHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var appointment models.Appointment
	if err := c.ShouldBindJSON(&appointment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.appointmentService.Update(uint(id), &appointment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "appointment updated successfully"})
}

func (h *AppointmentHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	if err := h.appointmentService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "appointment deleted successfully"})
}
