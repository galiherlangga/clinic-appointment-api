package handlers

import (
	"net/http"
	"strconv"

	"github.com/galiherlangga/clinic-appointment/models"
	"github.com/galiherlangga/clinic-appointment/services"
	"github.com/gin-gonic/gin"
)

type ProviderHandler struct {
	providerService services.ProviderService
}

func NewProviderHandler(providerService services.ProviderService) *ProviderHandler {
	return &ProviderHandler{providerService: providerService}
}

func (h *ProviderHandler) FindAll(c *gin.Context) {
	var filter models.ProviderFilter
	filter.Name = c.Query("name")
	filter.Email = c.Query("email")
	filter.Specialty = c.Query("specialty")

	providers, err := h.providerService.GetAll(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, providers)
}

func (h *ProviderHandler) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	provider, err := h.providerService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "provider not found"})
		return
	}

	c.JSON(http.StatusOK, provider)
}

func (h *ProviderHandler) Create(c *gin.Context) {
	var provider models.Provider
	if err := c.ShouldBindJSON(&provider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.providerService.Create(&provider); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, provider)
}

func (h *ProviderHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var provider models.Provider
	if err := c.ShouldBindJSON(&provider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.providerService.Update(uint(id), &provider); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "provider updated successfully"})
}

func (h *ProviderHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	if err := h.providerService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "provider deleted successfully"})
}
