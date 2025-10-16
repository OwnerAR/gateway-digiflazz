package services

import (
	"fmt"
	"time"

	"gateway-digiflazz/internal/models"
	"gateway-digiflazz/pkg/digiflazz"

	"github.com/sirupsen/logrus"
)

// PascabayarService handles Pascabayar transaction operations
type PascabayarService struct {
	digiflazzClient *digiflazz.Client
	logger          *logrus.Logger
}

// NewPascabayarService creates a new Pascabayar service
func NewPascabayarService(client *digiflazz.Client, logger *logrus.Logger) *PascabayarService {
	return &PascabayarService{
		digiflazzClient: client,
		logger:          logger,
	}
}

// CheckBill checks the Pascabayar bill before payment
func (s *PascabayarService) CheckBill(req models.PascabayarCheckRequest) (*models.PascabayarCheckResponse, error) {
	s.logger.WithFields(logrus.Fields{
		"ref_id":      req.RefID,
		"customer_no": req.CustomerNo,
		"buyer_sku":   req.BuyerSKU,
	}).Info("Checking Pascabayar bill")

	// Validate request
	if err := s.validateCheckRequest(req); err != nil {
		s.logger.WithError(err).Error("Pascabayar check request validation failed")
		return nil, err
	}

	// Call Digiflazz API to check bill
	resp, err := s.digiflazzClient.CheckPascabayarBill(req)
	if err != nil {
		s.logger.WithError(err).Error("Digiflazz Pascabayar check API call failed")
		return nil, fmt.Errorf("failed to check bill: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"ref_id": resp.Data.RefID,
		"amount": resp.Data.Amount,
		"status": resp.Data.Status,
		"rc":     resp.Data.RC,
	}).Info("Pascabayar bill check completed")

	return resp, nil
}

// PayBill processes the Pascabayar bill payment
func (s *PascabayarService) PayBill(req models.PascabayarPayRequest) (*models.PascabayarPayResponse, error) {
	s.logger.WithFields(logrus.Fields{
		"ref_id":      req.RefID,
		"customer_no": req.CustomerNo,
		"buyer_sku":   req.BuyerSKU,
		"amount":      req.Amount,
	}).Info("Processing Pascabayar bill payment")

	// Validate request
	if err := s.validatePayRequest(req); err != nil {
		s.logger.WithError(err).Error("Pascabayar pay request validation failed")
		return nil, err
	}

	// Call Digiflazz API to pay bill
	resp, err := s.digiflazzClient.PayPascabayarBill(req)
	if err != nil {
		s.logger.WithError(err).Error("Digiflazz Pascabayar payment API call failed")
		return nil, fmt.Errorf("failed to pay bill: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"ref_id": resp.Data.RefID,
		"amount": resp.Data.Amount,
		"status": resp.Data.Status,
		"rc":     resp.Data.RC,
		"sn":     resp.Data.SN,
	}).Info("Pascabayar bill payment completed")

	return resp, nil
}

// validateCheckRequest validates the check bill request
func (s *PascabayarService) validateCheckRequest(req models.PascabayarCheckRequest) error {
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

// validatePayRequest validates the pay bill request
func (s *PascabayarService) validatePayRequest(req models.PascabayarPayRequest) error {
	if req.RefID == "" {
		return fmt.Errorf("ref_id is required")
	}
	if req.CustomerNo == "" {
		return fmt.Errorf("customer_no is required")
	}
	if req.BuyerSKU == "" {
		return fmt.Errorf("buyer_sku is required")
	}
	if req.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	return nil
}

// CreatePascabayarTransaction creates a new Pascabayar transaction record
func (s *PascabayarService) CreatePascabayarTransaction(tx *models.PascabayarTransaction) error {
	s.logger.WithField("ref_id", tx.RefID).Info("Creating Pascabayar transaction record")
	
	// TODO: Save to database
	tx.CreatedAt = time.Now()
	tx.UpdatedAt = time.Now()
	
	return nil
}

// UpdatePascabayarTransaction updates an existing Pascabayar transaction
func (s *PascabayarService) UpdatePascabayarTransaction(tx *models.PascabayarTransaction) error {
	s.logger.WithField("ref_id", tx.RefID).Info("Updating Pascabayar transaction record")
	
	// TODO: Update in database
	tx.UpdatedAt = time.Now()
	
	return nil
}

// GetPascabayarTransaction retrieves a Pascabayar transaction by ref_id
func (s *PascabayarService) GetPascabayarTransaction(refID string) (*models.PascabayarTransaction, error) {
	s.logger.WithField("ref_id", refID).Info("Retrieving Pascabayar transaction")
	
	// TODO: Get from database
	// For now, return a mock transaction
	return &models.PascabayarTransaction{
		RefID: refID,
		Status: "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
