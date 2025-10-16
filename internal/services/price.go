package services

import (
	"fmt"

	"gateway-digiflazz/internal/models"
	"gateway-digiflazz/pkg/digiflazz"

	"github.com/sirupsen/logrus"
)

// PriceService handles price operations
type PriceService struct {
	digiflazzClient *digiflazz.Client
	logger          *logrus.Logger
}

// NewPriceService creates a new price service
func NewPriceService(client *digiflazz.Client, logger *logrus.Logger) *PriceService {
	return &PriceService{
		digiflazzClient: client,
		logger:          logger,
	}
}

// GetPrices retrieves the price list
func (s *PriceService) GetPrices(priceType string) (*models.PriceResponse, error) {
	s.logger.WithField("type", priceType).Info("Retrieving price list")

	// Validate price type
	if priceType != "" && priceType != "prabayar" && priceType != "pascabayar" {
		return nil, fmt.Errorf("invalid price type: %s. Must be 'prabayar' or 'pascabayar'", priceType)
	}

	// Call Digiflazz API
	resp, err := s.digiflazzClient.GetPrices(priceType)
	if err != nil {
		s.logger.WithError(err).Error("Digiflazz price API call failed")
		return nil, fmt.Errorf("failed to get prices: %w", err)
	}

	s.logger.WithField("product_count", len(resp.Data)).Info("Price list retrieved successfully")
	return resp, nil
}

// GetProductByCode retrieves a specific product by code
func (s *PriceService) GetProductByCode(code string) (*models.Product, error) {
	s.logger.WithField("code", code).Info("Retrieving product by code")

	// Get all products first
	resp, err := s.GetPrices("")
	if err != nil {
		return nil, err
	}

	// Find product by code
	for _, product := range resp.Data {
		if product.Code == code {
			s.logger.WithField("code", code).Info("Product found")
			return &product, nil
		}
	}

	s.logger.WithField("code", code).Warn("Product not found")
	return nil, fmt.Errorf("product with code %s not found", code)
}

// GetProductsByCategory retrieves products by category
func (s *PriceService) GetProductsByCategory(category string) ([]models.Product, error) {
	s.logger.WithField("category", category).Info("Retrieving products by category")

	// Get all products first
	resp, err := s.GetPrices("")
	if err != nil {
		return nil, err
	}

	// Filter by category
	var products []models.Product
	for _, product := range resp.Data {
		if product.Category == category {
			products = append(products, product)
		}
	}

	s.logger.WithFields(logrus.Fields{
		"category": category,
		"count":    len(products),
	}).Info("Products retrieved by category")

	return products, nil
}
