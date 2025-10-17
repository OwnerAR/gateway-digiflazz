package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gateway-digiflazz/internal/models"
	"gateway-digiflazz/pkg/digiflazz"

	"github.com/sirupsen/logrus"
)

// PLNInquiryService handles PLN inquiry operations with caching
type PLNInquiryService struct {
	digiflazzClient *digiflazz.Client
	logger          *logrus.Logger
	cache           CacheInterface
	config          models.PLNInquiryConfig
	stats           *models.PLNInquiryStats
}

// CacheInterface defines the interface for caching operations
type CacheInterface interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	ClearAll(ctx context.Context) error
	DeleteExpired(ctx context.Context) error
	GetStats(ctx context.Context) (map[string]interface{}, error)
}

// NewPLNInquiryService creates a new PLN inquiry service
func NewPLNInquiryService(client *digiflazz.Client, logger *logrus.Logger, cache CacheInterface) *PLNInquiryService {
	return &PLNInquiryService{
		digiflazzClient: client,
		logger:          logger,
		cache:           cache,
		config: models.PLNInquiryConfig{
			CacheEnabled:   true,
			CacheTTL:       0, // No expiration for static PLN data
			CacheKeyPrefix: "pln_inquiry:",
		},
		stats: &models.PLNInquiryStats{},
	}
}

// InquiryPLN performs PLN inquiry with caching strategy
func (s *PLNInquiryService) InquiryPLN(req models.PLNInquiryRequest, refID string) (*models.PLNInquiryResponse, error) {
	startTime := time.Now()
	s.stats.TotalRequests++

	s.logger.WithFields(logrus.Fields{
		"customer_no": req.CustomerNo,
		"cached":      s.config.CacheEnabled,
	}).Info("Processing PLN inquiry")

	// Check cache first if enabled
	if s.config.CacheEnabled {
		cached, err := s.getFromCache(req.CustomerNo)
		if err == nil && cached != nil {
			s.stats.CacheHits++
			s.logger.WithFields(logrus.Fields{
				"customer_no": req.CustomerNo,
				"ref_id":      refID,
			}).Info("PLN inquiry served from cache")
			
			// Update response time
			s.updateResponseTime(time.Since(startTime))
			
			// Build response with current ref_id
			response := s.buildResponseFromCache(cached)
			response.Data.RefID = refID // Use current ref_id
			response.Data.Message = cached.Message // Ensure message is included
			response.Message = "PLN inquiry completed successfully (cached)"
			response.Status = 1
			return response, nil
		}
		s.stats.CacheMisses++
		s.logger.WithFields(logrus.Fields{
			"customer_no": req.CustomerNo,
			"ref_id":      refID,
		}).Info("PLN inquiry cache miss")
	}

	// Call Digiflazz API
	s.stats.APIRequests++
	resp, err := s.digiflazzClient.InquiryPLN(req)
	if err != nil {
		s.stats.ErrorCount++
		s.logger.WithError(err).Error("Digiflazz PLN inquiry API call failed")
		return nil, fmt.Errorf("failed to inquiry PLN: %w", err)
	}

	// Cache the response if successful and caching is enabled
	if s.config.CacheEnabled && resp.Data.RC == "00" {
		if err := s.setToCache(req.CustomerNo, refID, resp); err != nil {
			s.logger.WithError(err).Warn("Failed to cache PLN inquiry response")
		}
	}

	// Update response time
	s.updateResponseTime(time.Since(startTime))

	// Ensure response has proper structure for cache miss
	if resp.Data.RefID == "" {
		resp.Data.RefID = refID
	}
	if resp.Data.Message == "" && resp.Data.RC == "00" {
		resp.Data.Message = "Transaksi Sukses"
	}
	if resp.Message == "" {
		resp.Message = "PLN inquiry completed successfully"
	}
	if resp.Status == 0 && resp.Data.RC == "00" {
		resp.Status = 1
	}

	s.logger.WithFields(logrus.Fields{
		"customer_no": req.CustomerNo,
		"rc":          resp.Data.RC,
		"status":      resp.Data.Status,
		"ref_id":      resp.Data.RefID,
		"message":     resp.Data.Message,
		"cached":      s.config.CacheEnabled,
	}).Info("PLN inquiry completed")

	return resp, nil
}

// getFromCache retrieves PLN inquiry data from cache
func (s *PLNInquiryService) getFromCache(customerNo string) (*models.PLNInquiryCache, error) {
	ctx := context.Background()
	key := s.getCacheKey(customerNo)
	
	cachedData, err := s.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var cache models.PLNInquiryCache
	if err := json.Unmarshal([]byte(cachedData), &cache); err != nil {
		return nil, err
	}

	// Check if cache is expired (only if ExpiresAt is set)
	if !cache.ExpiresAt.IsZero() && time.Now().After(cache.ExpiresAt) {
		// Delete expired cache
		s.cache.Delete(ctx, key)
		return nil, fmt.Errorf("cache expired")
	}

	return &cache, nil
}

// setToCache stores PLN inquiry data in cache
func (s *PLNInquiryService) setToCache(customerNo string, refID string, resp *models.PLNInquiryResponse) error {
	ctx := context.Background()
	key := s.getCacheKey(customerNo)

	cache := models.PLNInquiryCache{
		RefID:        refID,
		CustomerNo:   resp.Data.CustomerNo,
		MeterNo:      resp.Data.MeterNo,
		SubscriberID: resp.Data.SubscriberID,
		Name:         resp.Data.Name,
		SegmentPower: resp.Data.SegmentPower,
		Status:       resp.Data.Status,
		RC:           resp.Data.RC,
		Message:      resp.Data.Message,
		CachedAt:     time.Now(),
		ExpiresAt:    time.Time{}, // No expiration for static PLN data
	}

	cacheData, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	// Use 0 TTL for permanent cache (static PLN data)
	return s.cache.Set(ctx, key, string(cacheData), 0)
}

// getCacheKey generates cache key for customer number
func (s *PLNInquiryService) getCacheKey(customerNo string) string {
	return s.config.CacheKeyPrefix + customerNo
}

// buildResponseFromCache builds response from cached data
func (s *PLNInquiryService) buildResponseFromCache(cache *models.PLNInquiryCache) *models.PLNInquiryResponse {
	return &models.PLNInquiryResponse{
		Data: struct {
			Message       string `json:"message"`
			Status        string `json:"status"`
			RC            string `json:"rc"`
			RefID         string `json:"ref_id"`
			CustomerNo    string `json:"customer_no"`
			MeterNo       string `json:"meter_no,omitempty"`
			SubscriberID   string `json:"subscriber_id,omitempty"`
			Name          string `json:"name,omitempty"`
			SegmentPower  string `json:"segment_power,omitempty"`
		}{
			Message:      cache.Message,
			Status:       cache.Status,
			RC:           cache.RC,
			RefID:        cache.RefID,
			CustomerNo:   cache.CustomerNo,
			MeterNo:      cache.MeterNo,
			SubscriberID: cache.SubscriberID,
			Name:         cache.Name,
			SegmentPower: cache.SegmentPower,
		},
		Message: "PLN inquiry completed successfully (cached)",
		Status:  1,
	}
}

// updateResponseTime updates the average response time
func (s *PLNInquiryService) updateResponseTime(duration time.Duration) {
	if s.stats.AverageResponseTime == 0 {
		s.stats.AverageResponseTime = duration
	} else {
		s.stats.AverageResponseTime = (s.stats.AverageResponseTime + duration) / 2
	}
}

// GetStats returns PLN inquiry statistics
func (s *PLNInquiryService) GetStats() *models.PLNInquiryStats {
	return s.stats
}

// ClearCache clears PLN inquiry cache for a specific customer
func (s *PLNInquiryService) ClearCache(customerNo string) error {
	ctx := context.Background()
	key := s.getCacheKey(customerNo)
	
	s.logger.WithField("customer_no", customerNo).Info("Clearing PLN inquiry cache")
	return s.cache.Delete(ctx, key)
}

// ClearAllCache clears all PLN inquiry cache
func (s *PLNInquiryService) ClearAllCache() error {
	ctx := context.Background()
	s.logger.Info("Clearing all PLN inquiry cache")
	return s.cache.ClearAll(ctx)
}

// DeleteExpiredCache removes expired entries from cache
func (s *PLNInquiryService) DeleteExpiredCache() error {
	ctx := context.Background()
	s.logger.Info("Deleting expired PLN inquiry cache entries")
	return s.cache.DeleteExpired(ctx)
}

// GetCacheStats returns cache statistics
func (s *PLNInquiryService) GetCacheStats() (map[string]interface{}, error) {
	ctx := context.Background()
	return s.cache.GetStats(ctx)
}

// SetCacheConfig updates cache configuration
func (s *PLNInquiryService) SetCacheConfig(config models.PLNInquiryConfig) {
	s.config = config
	s.logger.WithFields(logrus.Fields{
		"cache_enabled": config.CacheEnabled,
		"cache_ttl":     config.CacheTTL,
	}).Info("PLN inquiry cache configuration updated")
}
