package handlers

import (
	"net/http"

	"gateway-digiflazz/internal/models"
	"gateway-digiflazz/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// PLNInquiryHandler handles PLN inquiry HTTP requests
type PLNInquiryHandler struct {
	plnInquiryService *services.PLNInquiryService
	logger            *logrus.Logger
}

// NewPLNInquiryHandler creates a new PLN inquiry handler
func NewPLNInquiryHandler(plnInquiryService *services.PLNInquiryService, logger *logrus.Logger) *PLNInquiryHandler {
	return &PLNInquiryHandler{
		plnInquiryService: plnInquiryService,
		logger:            logger,
	}
}

// InquiryPLN handles PLN inquiry requests
func (h *PLNInquiryHandler) InquiryPLN(c *gin.Context) {
	var req models.PLNInquiryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind PLN inquiry request")
		c.JSON(http.StatusBadRequest, models.PLNInquiryError{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Validate required fields
	if req.CustomerNo == "" {
		c.JSON(http.StatusBadRequest, models.PLNInquiryError{
			Code:    "MISSING_CUSTOMER_NO",
			Message: "customer_no is required",
		})
		return
	}

	// Perform PLN inquiry
	resp, err := h.plnInquiryService.InquiryPLN(req)
	if err != nil {
		h.logger.WithError(err).Error("PLN inquiry failed")
		c.JSON(http.StatusInternalServerError, models.PLNInquiryError{
			Code:    "INQUIRY_FAILED",
			Message: "Failed to perform PLN inquiry",
			Details: err.Error(),
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
	})
}

// GetStats handles PLN inquiry statistics requests
func (h *PLNInquiryHandler) GetStats(c *gin.Context) {
	stats := h.plnInquiryService.GetStats()
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// ClearCache handles cache clearing requests
func (h *PLNInquiryHandler) ClearCache(c *gin.Context) {
	customerNo := c.Param("customer_no")
	if customerNo == "" {
		c.JSON(http.StatusBadRequest, models.PLNInquiryError{
			Code:    "MISSING_CUSTOMER_NO",
			Message: "customer_no parameter is required",
		})
		return
	}

	// Clear cache for specific customer
	if err := h.plnInquiryService.ClearCache(customerNo); err != nil {
		h.logger.WithError(err).Error("Failed to clear PLN inquiry cache")
		c.JSON(http.StatusInternalServerError, models.PLNInquiryError{
			Code:    "CACHE_CLEAR_FAILED",
			Message: "Failed to clear cache",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cache cleared successfully",
	})
}

// ClearAllCache handles clearing all PLN inquiry cache
func (h *PLNInquiryHandler) ClearAllCache(c *gin.Context) {
	if err := h.plnInquiryService.ClearAllCache(); err != nil {
		h.logger.WithError(err).Error("Failed to clear all PLN inquiry cache")
		c.JSON(http.StatusInternalServerError, models.PLNInquiryError{
			Code:    "CACHE_CLEAR_FAILED",
			Message: "Failed to clear all cache",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "All cache cleared successfully",
	})
}

// UpdateCacheConfig handles cache configuration updates
func (h *PLNInquiryHandler) UpdateCacheConfig(c *gin.Context) {
	var config models.PLNInquiryConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		h.logger.WithError(err).Error("Failed to bind cache config request")
		c.JSON(http.StatusBadRequest, models.PLNInquiryError{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	h.plnInquiryService.SetCacheConfig(config)
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cache configuration updated successfully",
		"data":    config,
	})
}
