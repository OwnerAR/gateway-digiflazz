package services

import (
	"fmt"

	"gateway-digiflazz/internal/models"
	"gateway-digiflazz/pkg/digiflazz"

	"github.com/sirupsen/logrus"
)

// BalanceService handles balance operations
type BalanceService struct {
	digiflazzClient *digiflazz.Client
	logger          *logrus.Logger
}

// NewBalanceService creates a new balance service
func NewBalanceService(client *digiflazz.Client, logger *logrus.Logger) *BalanceService {
	return &BalanceService{
		digiflazzClient: client,
		logger:          logger,
	}
}

// GetBalance retrieves the current balance
func (s *BalanceService) GetBalance() (*models.BalanceResponse, error) {
	s.logger.Info("Retrieving account balance")

	// Call Digiflazz API
	resp, err := s.digiflazzClient.CheckBalance()
	if err != nil {
		s.logger.WithError(err).Error("Digiflazz balance API call failed")
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	s.logger.WithField("balance", resp.Data.Deposit).Info("Balance retrieved successfully")
	return resp, nil
}
