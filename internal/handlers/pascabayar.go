package handlers

import (
	"net/http"

	"gateway-digiflazz/internal/models"
	"gateway-digiflazz/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// PascabayarHandler handles Pascabayar HTTP requests
type PascabayarHandler struct {
	pascabayarService *services.PascabayarService
	logger            *logrus.Logger
}

// NewPascabayarHandler creates a new Pascabayar handler
func NewPascabayarHandler(pascabayarService *services.PascabayarService, logger *logrus.Logger) *PascabayarHandler {
	return &PascabayarHandler{
		pascabayarService: pascabayarService,
		logger:            logger,
	}
}

// CheckBill handles Pascabayar bill check requests
func (h *PascabayarHandler) CheckBill(c *gin.Context) {
	var req models.PascabayarCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind Pascabayar check request")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Check bill
	resp, err := h.pascabayarService.CheckBill(req)
	if err != nil {
		h.logger.WithError(err).Error("Pascabayar bill check failed")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    "BILL_CHECK_FAILED",
			Message: "Failed to check bill",
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

// PayBill handles Pascabayar bill payment requests
func (h *PascabayarHandler) PayBill(c *gin.Context) {
	var req models.PascabayarPayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind Pascabayar pay request")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Pay bill
	resp, err := h.pascabayarService.PayBill(req)
	if err != nil {
		h.logger.WithError(err).Error("Pascabayar bill payment failed")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    "BILL_PAYMENT_FAILED",
			Message: "Failed to pay bill",
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

// GetTransaction retrieves a Pascabayar transaction
func (h *PascabayarHandler) GetTransaction(c *gin.Context) {
	refID := c.Param("ref_id")
	if refID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    "MISSING_REF_ID",
			Message: "ref_id parameter is required",
		})
		return
	}

	// Get transaction
	tx, err := h.pascabayarService.GetPascabayarTransaction(refID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get Pascabayar transaction")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    "TRANSACTION_NOT_FOUND",
			Message: "Failed to get transaction",
			Details: err.Error(),
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tx,
	})
}
