package handlers

import (
	"net/http"

	"gateway-digiflazz/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// BalanceHandler handles balance HTTP requests
type BalanceHandler struct {
	balanceService *services.BalanceService
	logger         *logrus.Logger
}

// NewBalanceHandler creates a new balance handler
func NewBalanceHandler(balanceService *services.BalanceService, logger *logrus.Logger) *BalanceHandler {
	return &BalanceHandler{
		balanceService: balanceService,
		logger:         logger,
	}
}

// GetBalance handles balance check requests
func (h *BalanceHandler) GetBalance(c *gin.Context) {
	// Get balance
	resp, err := h.balanceService.GetBalance()
	if err != nil {
		h.logger.WithError(err).Error("Balance retrieval failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "BALANCE_FAILED",
				"message": "Failed to retrieve balance",
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
