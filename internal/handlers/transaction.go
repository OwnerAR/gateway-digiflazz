package handlers

import (
	"net/http"

	"gateway-digiflazz/internal/models"
	"gateway-digiflazz/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// TransactionHandler handles transaction HTTP requests
type TransactionHandler struct {
	transactionService *services.TransactionService
	logger             *logrus.Logger
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(transactionService *services.TransactionService, logger *logrus.Logger) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		logger:             logger,
	}
}

// Topup handles topup transaction requests
func (h *TransactionHandler) Topup(c *gin.Context) {
	var req models.TopupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind topup request")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Process topup
	resp, err := h.transactionService.Topup(req)
	if err != nil {
		h.logger.WithError(err).Error("Topup processing failed")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    "TOPUP_FAILED",
			Message: "Failed to process topup",
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

// Pay handles payment transaction requests
func (h *TransactionHandler) Pay(c *gin.Context) {
	var req models.PayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind payment request")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Process payment
	resp, err := h.transactionService.Pay(req)
	if err != nil {
		h.logger.WithError(err).Error("Payment processing failed")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    "PAYMENT_FAILED",
			Message: "Failed to process payment",
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

// GetStatus handles transaction status requests
func (h *TransactionHandler) GetStatus(c *gin.Context) {
	refID := c.Param("ref_id")
	if refID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    "MISSING_REF_ID",
			Message: "ref_id parameter is required",
		})
		return
	}

	// Get transaction status
	resp, err := h.transactionService.GetStatus(refID)
	if err != nil {
		h.logger.WithError(err).Error("Status check failed")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    "STATUS_CHECK_FAILED",
			Message: "Failed to check transaction status",
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

// Webhook handles webhook requests from Digiflazz
func (h *TransactionHandler) Webhook(c *gin.Context) {
	var webhook models.WebhookRequest
	if err := c.ShouldBindJSON(&webhook); err != nil {
		h.logger.WithError(err).Error("Failed to bind webhook request")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    "INVALID_WEBHOOK",
			Message: "Invalid webhook format",
			Details: err.Error(),
		})
		return
	}

	// Process webhook
	if err := h.transactionService.ProcessWebhook(webhook); err != nil {
		h.logger.WithError(err).Error("Webhook processing failed")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    "WEBHOOK_FAILED",
			Message: "Failed to process webhook",
			Details: err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Webhook processed successfully",
	})
}
