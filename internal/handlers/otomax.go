package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"gateway-digiflazz/internal/middleware"
	"gateway-digiflazz/internal/models"
	"gateway-digiflazz/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// OtomaxHandler handles Otomax HTTP requests
type OtomaxHandler struct {
	otomaxService    *services.OtomaxService
	plnInquiryService *services.PLNInquiryService
	logger           *logrus.Logger
}

// NewOtomaxHandler creates a new Otomax handler
func NewOtomaxHandler(otomaxService *services.OtomaxService, plnInquiryService *services.PLNInquiryService, logger *logrus.Logger) *OtomaxHandler {
	return &OtomaxHandler{
		otomaxService:    otomaxService,
		plnInquiryService: plnInquiryService,
		logger:           logger,
	}
}

// ProcessTransaction handles transaction requests from Otomax via GET with query parameters
func (h *OtomaxHandler) ProcessTransaction(c *gin.Context) {
	var req models.OtomaxTransactionRequest
	
	// Bind query parameters to struct
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind Otomax transaction request")
		middleware.ErrorResponse(c, http.StatusBadRequest, 
			"INVALID_REQUEST", 
			"Invalid request parameters", 
			err.Error())
		return
	}

	// Validate required fields
	if req.RefID == "" || req.CustomerNo == "" || req.BuyerSKU == "" {
		middleware.ErrorResponse(c, http.StatusBadRequest, 
			"MISSING_PARAMETERS", 
			"Missing required parameters: ref_id, customer_no, buyer_sku", 
			"")
		return
	}

	// Process transaction
	resp, err := h.otomaxService.ProcessTransaction(req)
	if err != nil {
		h.logger.WithError(err).Error("Otomax transaction processing failed")
		middleware.ErrorResponse(c, http.StatusInternalServerError, 
			"TRANSACTION_FAILED", 
			"Failed to process transaction", 
			err.Error())
		return
	}

	// Return response with consistent formatting
	middleware.SuccessResponse(c, resp, "Transaction processed successfully")
}

// CheckStatus handles status check requests from Otomax via GET with query parameters
func (h *OtomaxHandler) CheckStatus(c *gin.Context) {
	var req models.OtomaxStatusRequest
	
	// Bind query parameters to struct
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind Otomax status request")
		c.JSON(http.StatusBadRequest, models.OtomaxError{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request parameters",
			Details: err.Error(),
		})
		return
	}

	// Validate required fields
	if req.RefID == "" {
		c.JSON(http.StatusBadRequest, models.OtomaxError{
			Code:    "MISSING_PARAMETERS",
			Message: "Missing required parameter: ref_id",
		})
		return
	}

	// Check status
	resp, err := h.otomaxService.CheckStatus(req)
	if err != nil {
		h.logger.WithError(err).Error("Otomax status check failed")
		c.JSON(http.StatusInternalServerError, models.OtomaxError{
			Code:    "STATUS_CHECK_FAILED",
			Message: "Failed to check transaction status",
			Details: err.Error(),
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, resp)
}

// ProcessCallback handles callback requests from Digiflazz for Otomax transactions
func (h *OtomaxHandler) ProcessCallback(c *gin.Context) {
	var callback models.OtomaxCallback
	
	// Bind JSON body to struct
	if err := c.ShouldBindJSON(&callback); err != nil {
		h.logger.WithError(err).Error("Failed to bind Otomax callback")
		c.JSON(http.StatusBadRequest, models.OtomaxError{
			Code:    "INVALID_CALLBACK",
			Message: "Invalid callback format",
			Details: err.Error(),
		})
		return
	}

	// Process callback
	if err := h.otomaxService.ProcessCallback(callback); err != nil {
		h.logger.WithError(err).Error("Otomax callback processing failed")
		c.JSON(http.StatusInternalServerError, models.OtomaxError{
			Code:    "CALLBACK_FAILED",
			Message: "Failed to process callback",
			Details: err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Callback processed successfully",
	})
}

// GetTransactionHistory handles transaction history requests
func (h *OtomaxHandler) GetTransactionHistory(c *gin.Context) {
	// TODO: Implement transaction history retrieval
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Transaction history endpoint - to be implemented",
		"data":    []interface{}{},
	})
}

// GetProductList handles product list requests
func (h *OtomaxHandler) GetProductList(c *gin.Context) {
	// TODO: Implement product list retrieval
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Product list endpoint - to be implemented",
		"data":    []interface{}{},
	})
}

// CheckPascabayarBill handles Pascabayar bill check requests from Otomax
func (h *OtomaxHandler) CheckPascabayarBill(c *gin.Context) {
	var req models.OtomaxPascabayarCheckRequest
	
	// Bind query parameters to struct
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind Otomax Pascabayar check request")
		c.JSON(http.StatusBadRequest, models.OtomaxError{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request parameters",
			Details: err.Error(),
		})
		return
	}

	// Validate required fields
	if req.RefID == "" || req.CustomerNo == "" || req.BuyerSKU == "" {
		c.JSON(http.StatusBadRequest, models.OtomaxError{
			Code:    "MISSING_PARAMETERS",
			Message: "Missing required parameters: ref_id, customer_no, buyer_sku",
		})
		return
	}

	// TODO: Implement Pascabayar bill check logic
	// For now, return a mock response
	resp := models.OtomaxPascabayarCheckResponse{
		RefID:      req.RefID,
		CustomerNo: req.CustomerNo,
		BuyerSKU:   req.BuyerSKU,
		Amount:     50000,
		AdminFee:   2500,
		Total:      52500,
		Status:     "success",
		Message:    "Bill check successful",
		RC:         "00",
		BillDetails: models.BillDetails{
			CustomerName: "John Doe",
			BillPeriod:   "2023-12",
			DueDate:      "2023-12-31",
			BillAmount:   50000,
			AdminFee:     2500,
			TotalAmount:  52500,
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Sign:      h.generateResponseSignature(req.RefID, "success"),
	}

	// Return response
	c.JSON(http.StatusOK, resp)
}

// PayPascabayarBill handles Pascabayar bill payment requests from Otomax
func (h *OtomaxHandler) PayPascabayarBill(c *gin.Context) {
	var req models.OtomaxPascabayarPayRequest
	
	// Bind query parameters to struct
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind Otomax Pascabayar pay request")
		c.JSON(http.StatusBadRequest, models.OtomaxError{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request parameters",
			Details: err.Error(),
		})
		return
	}

	// Validate required fields
	if req.RefID == "" || req.CustomerNo == "" || req.BuyerSKU == "" {
		c.JSON(http.StatusBadRequest, models.OtomaxError{
			Code:    "MISSING_PARAMETERS",
			Message: "Missing required parameters: ref_id, customer_no, buyer_sku",
		})
		return
	}

	// TODO: Implement Pascabayar bill payment logic
	// For now, return a mock response
	resp := models.OtomaxPascabayarPayResponse{
		RefID:      req.RefID,
		CustomerNo: req.CustomerNo,
		BuyerSKU:   req.BuyerSKU,
		Amount:     req.Amount,
		AdminFee:   2500,
		Total:      req.Amount + 2500,
		Status:     "success",
		Message:    "Bill payment successful",
		RC:         "00",
		SN:         "SN123456789",
		BillDetails: models.BillDetails{
			CustomerName: "John Doe",
			BillPeriod:   "2023-12",
			DueDate:      "2023-12-31",
			BillAmount:   req.Amount,
			AdminFee:     2500,
			TotalAmount:  req.Amount + 2500,
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Sign:      h.generateResponseSignature(req.RefID, "success"),
	}

	// Return response
	c.JSON(http.StatusOK, resp)
}

// InquiryPLN handles PLN inquiry requests from Otomax
func (h *OtomaxHandler) InquiryPLN(c *gin.Context) {
	var req models.OtomaxPLNInquiryRequest
	
	// Bind query parameters to struct
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind PLN inquiry request")
		c.JSON(http.StatusBadRequest, models.OtomaxError{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request parameters",
			Details: err.Error(),
		})
		return
	}

	// Validate required fields
	if req.RefID == "" || req.CustomerNo == "" {
		c.JSON(http.StatusBadRequest, models.OtomaxError{
			Code:    "MISSING_PARAMETERS",
			Message: "Missing required parameters: ref_id, customer_no",
		})
		return
	}

	// Validate customer number format (PLN customer numbers are typically 11-12 digits)
	if len(req.CustomerNo) < 10 || len(req.CustomerNo) > 15 {
		c.JSON(http.StatusBadRequest, models.OtomaxError{
			Code:    "INVALID_CUSTOMER_NO",
			Message: "Invalid customer number format. PLN customer numbers should be 10-15 digits",
		})
		return
	}

	// Log request details for debugging
	h.logger.WithFields(logrus.Fields{
		"ref_id":      req.RefID,
		"customer_no": req.CustomerNo,
		"timestamp":   req.Timestamp,
	}).Info("Processing PLN inquiry request")

	// Use PLN inquiry service with cache strategy
	plnReq := models.PLNInquiryRequest{
		RefID:      req.RefID,
		CustomerNo: req.CustomerNo,
	}
	
	resp, err := h.plnInquiryService.InquiryPLN(plnReq)
	if err != nil {
		h.logger.WithError(err).WithFields(logrus.Fields{
			"ref_id":      req.RefID,
			"customer_no": req.CustomerNo,
		}).Error("PLN inquiry failed")
		
		// Check for specific error types
		if strings.Contains(err.Error(), "customer may not exist") {
			middleware.ErrorResponse(c, http.StatusNotFound, 
				"CUSTOMER_NOT_FOUND", 
				"Customer number not found in PLN system", 
				fmt.Sprintf("Customer number %s does not exist or is invalid", req.CustomerNo))
		} else if strings.Contains(err.Error(), "IP whitelist error") {
			middleware.ErrorResponse(c, http.StatusForbidden, 
				"IP_NOT_WHITELISTED", 
				"Server IP not registered in Digiflazz whitelist", 
				"Please contact Digiflazz to add your server IP to the whitelist")
		} else {
			middleware.ErrorResponse(c, http.StatusInternalServerError, 
				"PLN_INQUIRY_FAILED", 
				"Failed to perform PLN inquiry", 
				err.Error())
		}
		return
	}

	// Log successful inquiry
	h.logger.WithFields(logrus.Fields{
		"ref_id":      req.RefID,
		"customer_no": req.CustomerNo,
		"rc":          resp.Data.RC,
		"status":      resp.Data.Status,
		"meter_no":    resp.Data.MeterNo,
		"name":        resp.Data.Name,
	}).Info("PLN inquiry completed successfully")

	// Return response with consistent formatting
	middleware.SuccessResponse(c, resp, "PLN inquiry completed successfully")
}

// GetPLNStats handles PLN inquiry statistics requests from Otomax
func (h *OtomaxHandler) GetPLNStats(c *gin.Context) {
	stats := h.plnInquiryService.GetStats()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "PLN inquiry statistics",
		"data":    stats,
	})
}

// ClearPLNCache handles PLN cache clearing requests from Otomax
func (h *OtomaxHandler) ClearPLNCache(c *gin.Context) {
	customerNo := c.Param("customer_no")
	if customerNo == "" {
		c.JSON(http.StatusBadRequest, models.OtomaxError{
			Code:    "MISSING_PARAMETERS",
			Message: "Missing required parameter: customer_no",
		})
		return
	}

	err := h.plnInquiryService.ClearCache(customerNo)
	if err != nil {
		h.logger.WithError(err).Error("Failed to clear PLN cache")
		c.JSON(http.StatusInternalServerError, models.OtomaxError{
			Code:    "CACHE_CLEAR_FAILED",
			Message: "Failed to clear PLN cache",
			Details: err.Error(),
		})
		return
	}

	h.logger.WithField("customer_no", customerNo).Info("PLN cache cleared for customer")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "PLN cache cleared for customer: " + customerNo,
	})
}

// ClearAllPLNCache handles clearing all PLN cache requests from Otomax
func (h *OtomaxHandler) ClearAllPLNCache(c *gin.Context) {
	err := h.plnInquiryService.ClearAllCache()
	if err != nil {
		h.logger.WithError(err).Error("Failed to clear all PLN cache")
		c.JSON(http.StatusInternalServerError, models.OtomaxError{
			Code:    "CACHE_CLEAR_ALL_FAILED",
			Message: "Failed to clear all PLN cache",
			Details: err.Error(),
		})
		return
	}

	h.logger.Info("All PLN cache cleared")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "All PLN cache cleared",
	})
}

// UpdatePLNCacheConfig handles PLN cache configuration updates from Otomax
func (h *OtomaxHandler) UpdatePLNCacheConfig(c *gin.Context) {
	var config models.PLNInquiryConfig
	
	// Bind JSON body to struct
	if err := c.ShouldBindJSON(&config); err != nil {
		h.logger.WithError(err).Error("Failed to bind PLN cache config")
		c.JSON(http.StatusBadRequest, models.OtomaxError{
			Code:    "INVALID_REQUEST",
			Message: "Invalid configuration format",
			Details: err.Error(),
		})
		return
	}

	// TODO: Implement PLN cache configuration update logic
	h.logger.WithFields(map[string]interface{}{
		"cache_enabled": config.CacheEnabled,
		"cache_ttl":     config.CacheTTL,
	}).Info("Updating PLN cache configuration")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "PLN cache configuration updated",
		"data":    config,
	})
}

// generateResponseSignature generates signature for response
func (h *OtomaxHandler) generateResponseSignature(refID, status string) string {
	// TODO: Implement proper signature generation
	return "mock_signature_" + refID + "_" + status
}
