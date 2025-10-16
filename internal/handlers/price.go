package handlers

import (
	"net/http"

	"gateway-digiflazz/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// PriceHandler handles price HTTP requests
type PriceHandler struct {
	priceService *services.PriceService
	logger       *logrus.Logger
}

// NewPriceHandler creates a new price handler
func NewPriceHandler(priceService *services.PriceService, logger *logrus.Logger) *PriceHandler {
	return &PriceHandler{
		priceService: priceService,
		logger:       logger,
	}
}

// GetPrices handles price list requests
func (h *PriceHandler) GetPrices(c *gin.Context) {
	// Get price type from query parameter
	priceType := c.Query("type")

	// Get prices
	resp, err := h.priceService.GetPrices(priceType)
	if err != nil {
		h.logger.WithError(err).Error("Price list retrieval failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "PRICES_FAILED",
				"message": "Failed to retrieve price list",
				"details": err.Error(),
			},
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
	})
}

// GetProductByCode handles product lookup by code
func (h *PriceHandler) GetProductByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_CODE",
				"message": "Product code is required",
			},
		})
		return
	}

	// Get product
	product, err := h.priceService.GetProductByCode(code)
	if err != nil {
		h.logger.WithError(err).Error("Product lookup failed")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "PRODUCT_NOT_FOUND",
				"message": "Product not found",
				"details": err.Error(),
			},
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    product,
	})
}

// GetProductsByCategory handles product lookup by category
func (h *PriceHandler) GetProductsByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_CATEGORY",
				"message": "Category is required",
			},
		})
		return
	}

	// Get products
	products, err := h.priceService.GetProductsByCategory(category)
	if err != nil {
		h.logger.WithError(err).Error("Products lookup failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "PRODUCTS_FAILED",
				"message": "Failed to retrieve products",
				"details": err.Error(),
			},
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    products,
	})
}
