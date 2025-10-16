package services

import (
	"fmt"
	"time"

	"gateway-digiflazz/internal/models"
	"gateway-digiflazz/pkg/digiflazz"

	"github.com/sirupsen/logrus"
)

// TransactionService handles transaction operations
type TransactionService struct {
	digiflazzClient *digiflazz.Client
	logger          *logrus.Logger
}

// NewTransactionService creates a new transaction service
func NewTransactionService(client *digiflazz.Client, logger *logrus.Logger) *TransactionService {
	return &TransactionService{
		digiflazzClient: client,
		logger:          logger,
	}
}

// Topup performs a topup transaction
func (s *TransactionService) Topup(req models.TopupRequest) (*models.TopupResponse, error) {
	s.logger.WithFields(logrus.Fields{
		"ref_id":      req.RefID,
		"customer_no": req.CustomerNo,
		"buyer_sku":   req.BuyerSKU,
	}).Info("Processing topup transaction")

	// Validate request
	if err := s.validateTopupRequest(req); err != nil {
		s.logger.WithError(err).Error("Topup request validation failed")
		return nil, err
	}

	// Call Digiflazz API
	resp, err := s.digiflazzClient.Topup(req)
	if err != nil {
		s.logger.WithError(err).Error("Digiflazz topup API call failed")
		return nil, fmt.Errorf("failed to process topup: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"ref_id": resp.Data.RefID,
		"status": resp.Data.Status,
		"rc":     resp.Data.RC,
	}).Info("Topup transaction completed")

	return resp, nil
}

// Pay performs a payment transaction
func (s *TransactionService) Pay(req models.PayRequest) (*models.PayResponse, error) {
	s.logger.WithFields(logrus.Fields{
		"ref_id":      req.RefID,
		"customer_no": req.CustomerNo,
		"buyer_sku":   req.BuyerSKU,
	}).Info("Processing payment transaction")

	// Validate request
	if err := s.validatePayRequest(req); err != nil {
		s.logger.WithError(err).Error("Payment request validation failed")
		return nil, err
	}

	// Call Digiflazz API
	resp, err := s.digiflazzClient.Pay(req)
	if err != nil {
		s.logger.WithError(err).Error("Digiflazz payment API call failed")
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"ref_id": resp.Data.RefID,
		"status": resp.Data.Status,
		"rc":     resp.Data.RC,
	}).Info("Payment transaction completed")

	return resp, nil
}

// GetStatus checks the transaction status
func (s *TransactionService) GetStatus(refID string) (*models.StatusResponse, error) {
	s.logger.WithField("ref_id", refID).Info("Checking transaction status")

	// Validate refID
	if refID == "" {
		return nil, fmt.Errorf("ref_id is required")
	}

	// Call Digiflazz API
	resp, err := s.digiflazzClient.CheckStatus(refID)
	if err != nil {
		s.logger.WithError(err).Error("Digiflazz status check API call failed")
		return nil, fmt.Errorf("failed to check status: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"ref_id": resp.Data.RefID,
		"status": resp.Data.Status,
		"rc":     resp.Data.RC,
	}).Info("Transaction status retrieved")

	return resp, nil
}

// ProcessWebhook processes incoming webhook from Digiflazz
func (s *TransactionService) ProcessWebhook(webhook models.WebhookRequest) error {
	s.logger.WithFields(logrus.Fields{
		"ref_id": webhook.RefID,
		"status": webhook.Status,
		"rc":     webhook.RC,
	}).Info("Processing webhook")

	// Validate webhook signature
	if !s.digiflazzClient.ValidateWebhook(webhook) {
		s.logger.Error("Invalid webhook signature")
		return fmt.Errorf("invalid webhook signature")
	}

	// TODO: Update transaction status in database
	// TODO: Send notification to user
	// TODO: Update internal systems

	s.logger.WithField("ref_id", webhook.RefID).Info("Webhook processed successfully")
	return nil
}

// validateTopupRequest validates topup request
func (s *TransactionService) validateTopupRequest(req models.TopupRequest) error {
	if req.RefID == "" {
		return fmt.Errorf("ref_id is required")
	}
	if req.CustomerNo == "" {
		return fmt.Errorf("customer_no is required")
	}
	if req.BuyerSKU == "" {
		return fmt.Errorf("buyer_sku is required")
	}
	return nil
}

// validatePayRequest validates payment request
func (s *TransactionService) validatePayRequest(req models.PayRequest) error {
	if req.RefID == "" {
		return fmt.Errorf("ref_id is required")
	}
	if req.CustomerNo == "" {
		return fmt.Errorf("customer_no is required")
	}
	if req.BuyerSKU == "" {
		return fmt.Errorf("buyer_sku is required")
	}
	return nil
}

// CreateTransaction creates a new transaction record
func (s *TransactionService) CreateTransaction(tx *models.Transaction) error {
	s.logger.WithField("ref_id", tx.RefID).Info("Creating transaction record")
	
	// TODO: Save to database
	tx.CreatedAt = time.Now()
	tx.UpdatedAt = time.Now()
	
	return nil
}

// UpdateTransaction updates an existing transaction
func (s *TransactionService) UpdateTransaction(tx *models.Transaction) error {
	s.logger.WithField("ref_id", tx.RefID).Info("Updating transaction record")
	
	// TODO: Update in database
	tx.UpdatedAt = time.Now()
	
	return nil
}
