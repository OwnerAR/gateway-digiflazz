package services

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"

	"gateway-digiflazz/internal/models"
	"gateway-digiflazz/pkg/digiflazz"

	"github.com/sirupsen/logrus"
)

// OtomaxService handles Otomax transaction operations
type OtomaxService struct {
	digiflazzClient *digiflazz.Client
	logger          *logrus.Logger
	secretKey       string
}

// NewOtomaxService creates a new Otomax service
func NewOtomaxService(client *digiflazz.Client, logger *logrus.Logger, secretKey string) *OtomaxService {
	return &OtomaxService{
		digiflazzClient: client,
		logger:          logger,
		secretKey:       secretKey,
	}
}

// ProcessTransaction processes a transaction from Otomax
func (s *OtomaxService) ProcessTransaction(req models.OtomaxTransactionRequest) (*models.OtomaxTransactionResponse, error) {
	s.logger.WithFields(logrus.Fields{
		"ref_id":      req.RefID,
		"customer_no": req.CustomerNo,
		"buyer_sku":   req.BuyerSKU,
		"type":        req.Type,
	}).Info("Processing Otomax transaction")

	// Note: Signature validation removed - Otomax requests do not require signature validation

	// Parse amount
	amount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		s.logger.WithError(err).Error("Invalid amount format")
		return nil, fmt.Errorf("invalid amount format")
	}

	// Create transaction record
	transaction := &models.OtomaxTransaction{
		RefID:      req.RefID,
		CustomerNo: req.CustomerNo,
		BuyerSKU:   req.BuyerSKU,
		Amount:     amount,
		Type:       req.Type,
		Status:     "pending",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Process based on transaction type
	var response *models.OtomaxTransactionResponse
	if req.Type == "prabayar" {
		response, err = s.processPrabayarTransaction(transaction)
	} else if req.Type == "pascabayar" {
		response, err = s.processPascabayarTransaction(transaction)
	} else {
		return nil, fmt.Errorf("invalid transaction type: %s", req.Type)
	}

	if err != nil {
		s.logger.WithError(err).Error("Transaction processing failed")
		transaction.Status = "failed"
		transaction.Message = err.Error()
		// TODO: Update transaction in database
		return nil, err
	}

	// TODO: Save transaction to database
	s.logger.WithField("ref_id", req.RefID).Info("Otomax transaction processed successfully")
	return response, nil
}

// CheckStatus checks the status of an Otomax transaction
func (s *OtomaxService) CheckStatus(req models.OtomaxStatusRequest) (*models.OtomaxStatusResponse, error) {
	s.logger.WithField("ref_id", req.RefID).Info("Checking Otomax transaction status")

	// Note: Signature validation removed - Otomax requests do not require signature validation

	// TODO: Get transaction from database
	// For now, return a mock response
	response := &models.OtomaxStatusResponse{
		RefID:     req.RefID,
		Status:    "success",
		Message:   "Transaction completed",
		RC:        "00",
		Timestamp: time.Now().Format(time.RFC3339),
		Sign:      s.generateResponseSignature(req.RefID, "success"),
	}

	s.logger.WithField("ref_id", req.RefID).Info("Otomax transaction status retrieved")
	return response, nil
}

// ProcessCallback processes callback from Digiflazz for Otomax transactions
func (s *OtomaxService) ProcessCallback(callback models.OtomaxCallback) error {
	s.logger.WithFields(logrus.Fields{
		"ref_id": callback.RefID,
		"status": callback.Status,
		"rc":     callback.RC,
	}).Info("Processing Otomax callback")

	// Validate callback signature
	if !s.validateCallbackSignature(callback) {
		s.logger.Error("Invalid callback signature")
		return fmt.Errorf("invalid callback signature")
	}

	// TODO: Update transaction status in database
	// TODO: Notify Otomax about status update
	// TODO: Send notification to user

	s.logger.WithField("ref_id", callback.RefID).Info("Otomax callback processed successfully")
	return nil
}

// processPrabayarTransaction processes a prabayar transaction
func (s *OtomaxService) processPrabayarTransaction(transaction *models.OtomaxTransaction) (*models.OtomaxTransactionResponse, error) {
	// Create Digiflazz topup request
	digiflazzReq := models.TopupRequest{
		RefID:      transaction.RefID,
		CustomerNo: transaction.CustomerNo,
		BuyerSKU:   transaction.BuyerSKU,
	}

	// Call Digiflazz API
	digiflazzResp, err := s.digiflazzClient.Topup(digiflazzReq)
	if err != nil {
		return nil, fmt.Errorf("digiflazz topup failed: %w", err)
	}

	// Create response
	response := &models.OtomaxTransactionResponse{
		RefID:      transaction.RefID,
		CustomerNo: transaction.CustomerNo,
		BuyerSKU:   transaction.BuyerSKU,
		Amount:     transaction.Amount,
		Status:     s.mapDigiflazzStatus(digiflazzResp.Data.Status),
		Message:    digiflazzResp.Data.Message,
		RC:         digiflazzResp.Data.RC,
		SN:         digiflazzResp.Data.SN,
		Timestamp:  time.Now().Format(time.RFC3339),
		Sign:       s.generateResponseSignature(transaction.RefID, s.mapDigiflazzStatus(digiflazzResp.Data.Status)),
	}

	// Update transaction status
	transaction.Status = response.Status
	transaction.Message = response.Message
	transaction.RC = response.RC
	transaction.SN = response.SN
	transaction.DigiflazzRefID = digiflazzResp.Data.RefID

	return response, nil
}

// processPascabayarTransaction processes a pascabayar transaction
func (s *OtomaxService) processPascabayarTransaction(transaction *models.OtomaxTransaction) (*models.OtomaxTransactionResponse, error) {
	// For Pascabayar, we need to check the bill first
	// This is a two-step process: Check -> Pay
	
	// Step 1: Check the bill
	checkReq := models.PascabayarCheckRequest{
		RefID:      transaction.RefID,
		CustomerNo: transaction.CustomerNo,
		BuyerSKU:   transaction.BuyerSKU,
	}

	checkResp, err := s.digiflazzClient.CheckPascabayarBill(checkReq)
	if err != nil {
		return nil, fmt.Errorf("digiflazz bill check failed: %w", err)
	}

	// If check failed, return error
	if checkResp.Data.RC != "00" {
		return &models.OtomaxTransactionResponse{
			RefID:      transaction.RefID,
			CustomerNo: transaction.CustomerNo,
			BuyerSKU:   transaction.BuyerSKU,
			Amount:     transaction.Amount,
			Status:     "failed",
			Message:    checkResp.Data.Message,
			RC:         checkResp.Data.RC,
			Timestamp:  time.Now().Format(time.RFC3339),
			Sign:       s.generateResponseSignature(transaction.RefID, "failed"),
		}, nil
	}

	// Step 2: Pay the bill
	payReq := models.PascabayarPayRequest{
		RefID:      transaction.RefID,
		CustomerNo: transaction.CustomerNo,
		BuyerSKU:   transaction.BuyerSKU,
		Amount:     checkResp.Data.Amount, // Use amount from check response
	}

	payResp, err := s.digiflazzClient.PayPascabayarBill(payReq)
	if err != nil {
		return nil, fmt.Errorf("digiflazz bill payment failed: %w", err)
	}

	// Create response
	response := &models.OtomaxTransactionResponse{
		RefID:      transaction.RefID,
		CustomerNo: transaction.CustomerNo,
		BuyerSKU:   transaction.BuyerSKU,
		Amount:     payResp.Data.Amount,
		Status:     s.mapDigiflazzStatus(payResp.Data.Status),
		Message:    payResp.Data.Message,
		RC:         payResp.Data.RC,
		Timestamp:  time.Now().Format(time.RFC3339),
		Sign:       s.generateResponseSignature(transaction.RefID, s.mapDigiflazzStatus(payResp.Data.Status)),
	}

	// Update transaction status
	transaction.Status = response.Status
	transaction.Message = response.Message
	transaction.RC = response.RC
	transaction.DigiflazzRefID = payResp.Data.RefID

	return response, nil
}

// Note: Signature validation removed - Otomax requests do not require signature validation
// The gateway handles all Digiflazz API signatures internally

// validateCallbackSignature validates the callback signature
func (s *OtomaxService) validateCallbackSignature(callback models.OtomaxCallback) bool {
	expectedSign := s.generateCallbackSignature(callback.RefID, callback.Status)
	return callback.Sign == expectedSign
}

// generateRequestSignature generates signature for transaction request
func (s *OtomaxService) generateRequestSignature(refID, customerNo, buyerSKU string) string {
	data := fmt.Sprintf("%s%s%s%s", refID, customerNo, buyerSKU, s.secretKey)
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

// generateStatusSignature generates signature for status request
func (s *OtomaxService) generateStatusSignature(refID string) string {
	data := fmt.Sprintf("%s%s", refID, s.secretKey)
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

// generateResponseSignature generates signature for response
func (s *OtomaxService) generateResponseSignature(refID, status string) string {
	data := fmt.Sprintf("%s%s%s", refID, status, s.secretKey)
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

// generateCallbackSignature generates signature for callback
func (s *OtomaxService) generateCallbackSignature(refID, status string) string {
	data := fmt.Sprintf("%s%s%s", refID, status, s.secretKey)
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

// mapDigiflazzStatus maps Digiflazz status to Otomax status
func (s *OtomaxService) mapDigiflazzStatus(digiflazzStatus string) string {
	switch digiflazzStatus {
	case "success":
		return "success"
	case "pending":
		return "pending"
	case "failed":
		return "failed"
	default:
		return "failed"
	}
}
