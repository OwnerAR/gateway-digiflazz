package models

import "time"

// PLNInquiryRequest represents the request for PLN inquiry
type PLNInquiryRequest struct {
	Username   string `json:"username" binding:"required"`
	CustomerNo string `json:"customer_no" binding:"required"`
	RefID      string `json:"ref_id" binding:"required"`
	Sign       string `json:"sign" binding:"required"`
}

// OtomaxPLNInquiryRequest represents the request for PLN inquiry from Otomax
type OtomaxPLNInquiryRequest struct {
	RefID      string `form:"ref_id" json:"ref_id" binding:"required"`
	CustomerNo string `form:"customer_no" json:"customer_no" binding:"required"`
	Timestamp  string `form:"timestamp" json:"timestamp"`
}

// PLNInquiryResponse represents the response for PLN inquiry
type PLNInquiryResponse struct {
	Data struct {
		Message       string `json:"message"`
		Status        string `json:"status"`
		RC            string `json:"rc"`
		RefID         string `json:"ref_id"`
		CustomerNo    string `json:"customer_no"`
		MeterNo       string `json:"meter_no,omitempty"`
		SubscriberID   string `json:"subscriber_id,omitempty"`
		Name          string `json:"name,omitempty"`
		SegmentPower  string `json:"segment_power,omitempty"`
	} `json:"data"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// PLNInquiryCache represents cached PLN inquiry data
type PLNInquiryCache struct {
	RefID         string    `json:"ref_id"`
	CustomerNo    string    `json:"customer_no"`
	MeterNo       string    `json:"meter_no"`
	SubscriberID  string    `json:"subscriber_id"`
	Name          string    `json:"name"`
	SegmentPower  string    `json:"segment_power"`
	Status        string    `json:"status"`
	RC            string    `json:"rc"`
	Message       string    `json:"message"`
	CachedAt      time.Time `json:"cached_at"`
	ExpiresAt     time.Time `json:"expires_at"`
}

// PLNInquiryError represents an error response for PLN inquiry
type PLNInquiryError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// PLNInquiryConfig represents configuration for PLN inquiry caching
type PLNInquiryConfig struct {
	CacheEnabled   bool          `json:"cache_enabled"`
	CacheTTL       time.Duration `json:"cache_ttl"`
	CacheKeyPrefix string        `json:"cache_key_prefix"`
}

// PLNInquiryStats represents statistics for PLN inquiry
type PLNInquiryStats struct {
	TotalRequests    int64 `json:"total_requests"`
	CacheHits        int64 `json:"cache_hits"`
	CacheMisses      int64 `json:"cache_misses"`
	APIRequests      int64 `json:"api_requests"`
	ErrorCount       int64 `json:"error_count"`
	AverageResponseTime time.Duration `json:"average_response_time"`
}
