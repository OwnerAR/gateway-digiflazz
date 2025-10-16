package digiflazz

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gateway-digiflazz/internal/config"
	"gateway-digiflazz/internal/models"

	"github.com/sirupsen/logrus"
)

// Client represents the Digiflazz API client
type Client struct {
	config     config.DigiflazzConfig
	httpClient *http.Client
	baseURL    string
	logger     *logrus.Logger
}

// NewClient creates a new Digiflazz API client
func NewClient(cfg config.DigiflazzConfig, logger *logrus.Logger) *Client {
	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		baseURL: cfg.BaseURL,
		logger:  logger,
	}
}

// generateSign generates MD5 signature for Digiflazz API
func (c *Client) generateSign(username, apiKey, refID string) string {
	data := fmt.Sprintf("%s%s%s", username, apiKey, refID)
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

// CheckBalance checks the account balance
func (c *Client) CheckBalance() (*models.BalanceResponse, error) {
	req := models.BalanceRequest{
		DigiflazzRequest: models.DigiflazzRequest{
			Username: c.config.Username,
			APIKey:   c.config.APIKey,
			Sign:     c.generateSign(c.config.Username, c.config.APIKey, "deposit"),
		},
	}

	var resp models.BalanceResponse
	if err := c.makeRequest("/cek-saldo", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetPrices gets the price list
func (c *Client) GetPrices(priceType string) (*models.PriceResponse, error) {
	req := models.PriceRequest{
		DigiflazzRequest: models.DigiflazzRequest{
			Username: c.config.Username,
			APIKey:   c.config.APIKey,
			Sign:     c.generateSign(c.config.Username, c.config.APIKey, "pricelist"),
		},
		Type: priceType,
	}

	var resp models.PriceResponse
	if err := c.makeRequest("/daftar-harga", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Topup performs a topup transaction
func (c *Client) Topup(req models.TopupRequest) (*models.TopupResponse, error) {
	// Generate signature for topup
	req.Sign = c.generateSign(c.config.Username, c.config.APIKey, req.RefID)
	req.Username = c.config.Username
	req.APIKey = c.config.APIKey

	var resp models.TopupResponse
	if err := c.makeRequest("/topup", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Pay performs a payment transaction
func (c *Client) Pay(req models.PayRequest) (*models.PayResponse, error) {
	// Generate signature for payment
	req.Sign = c.generateSign(c.config.Username, c.config.APIKey, req.RefID)
	req.Username = c.config.Username
	req.APIKey = c.config.APIKey

	var resp models.PayResponse
	if err := c.makeRequest("/pascabayar", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// CheckStatus checks the transaction status
func (c *Client) CheckStatus(refID string) (*models.StatusResponse, error) {
	req := models.StatusRequest{
		DigiflazzRequest: models.DigiflazzRequest{
			Username: c.config.Username,
			APIKey:   c.config.APIKey,
			Sign:     c.generateSign(c.config.Username, c.config.APIKey, refID),
		},
		RefID: refID,
	}

	var resp models.StatusResponse
	if err := c.makeRequest("/cek-status", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// makeRequest makes an HTTP request to Digiflazz API
func (c *Client) makeRequest(endpoint string, req interface{}, resp interface{}) error {
	// Marshal request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := c.baseURL + endpoint
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", "Digiflazz-Gateway/1.0")

	// Log request details
	c.logger.WithFields(logrus.Fields{
		"url":          url,
		"method":       "POST",
		"endpoint":     endpoint,
		"payload":      string(jsonData),
		"payload_size": len(jsonData),
		"timeout":      c.httpClient.Timeout,
	}).Debug("Making request to Digiflazz API")
	
	// Log full JSON payload for debugging
	c.logger.WithField("json_payload", string(jsonData)).Info("Full JSON payload being sent to Digiflazz")
	
	// Log environment information
	c.logger.WithFields(logrus.Fields{
		"base_url":     c.baseURL,
		"username":     c.config.Username,
		"api_key_len":  len(c.config.APIKey),
		"timeout":      c.config.Timeout,
		"retry_attempts": c.config.RetryAttempts,
	}).Info("Digiflazz client configuration")

	// Make request with retry logic
	var lastErr error
	for attempt := 0; attempt < c.config.RetryAttempts; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * time.Second)
		}

		httpResp, err := c.httpClient.Do(httpReq)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			continue
		}
		defer httpResp.Body.Close()

		// Read response body
		body, err := io.ReadAll(httpResp.Body)
		if err != nil {
			lastErr = fmt.Errorf("failed to read response: %w", err)
			continue
		}

		// Log response details
		c.logger.WithFields(logrus.Fields{
			"status_code": httpResp.StatusCode,
			"response":    string(body),
			"endpoint":    endpoint,
		}).Debug("Received response from Digiflazz API")

		// Check HTTP status
		if httpResp.StatusCode != http.StatusOK {
			c.logger.WithFields(logrus.Fields{
				"status_code": httpResp.StatusCode,
				"response":    string(body),
				"endpoint":    endpoint,
			}).Error("Digiflazz API returned error status")
			lastErr = fmt.Errorf("HTTP error %d: %s", httpResp.StatusCode, string(body))
			continue
		}

		// Unmarshal response
		if err := json.Unmarshal(body, resp); err != nil {
			c.logger.WithError(err).Error("Failed to unmarshal Digiflazz API response")
			lastErr = fmt.Errorf("failed to unmarshal response: %w", err)
			continue
		}

		return nil
	}

	return lastErr
}

// CheckPascabayarBill checks the Pascabayar bill
func (c *Client) CheckPascabayarBill(req models.PascabayarCheckRequest) (*models.PascabayarCheckResponse, error) {
	// Generate signature for check request
	req.Sign = c.generateSign(c.config.Username, c.config.APIKey, req.RefID)
	req.Username = c.config.Username
	req.APIKey = c.config.APIKey

	var resp models.PascabayarCheckResponse
	if err := c.makeRequest("/pascabayar/check", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// PayPascabayarBill pays the Pascabayar bill
func (c *Client) PayPascabayarBill(req models.PascabayarPayRequest) (*models.PascabayarPayResponse, error) {
	// Generate signature for pay request
	req.Sign = c.generateSign(c.config.Username, c.config.APIKey, req.RefID)
	req.Username = c.config.Username
	req.APIKey = c.config.APIKey

	var resp models.PascabayarPayResponse
	if err := c.makeRequest("/pascabayar/pay", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// InquiryPLN performs PLN inquiry
func (c *Client) InquiryPLN(req models.PLNInquiryRequest) (*models.PLNInquiryResponse, error) {
	// Log request details
	c.logger.WithFields(logrus.Fields{
		"customer_no": req.CustomerNo,
		"endpoint":    "/inquiry-pln",
		"username":    c.config.Username,
	}).Info("Sending PLN inquiry request to Digiflazz API")

	// Generate signature for PLN inquiry
	req.Username = c.config.Username
	req.Sign = c.generatePLNInquirySign(c.config.Username, c.config.APIKey, req.CustomerNo)
	
	// Log signature generation details
	signatureInput := fmt.Sprintf("%s%s%s", c.config.Username, c.config.APIKey, req.CustomerNo)
	c.logger.WithFields(logrus.Fields{
		"username":        c.config.Username,
		"customer_no":     req.CustomerNo,
		"signature":       req.Sign,
		"signature_input": signatureInput,
		"signature_length": len(req.Sign),
		"api_key_length":  len(c.config.APIKey),
		"base_url":        c.baseURL,
	}).Info("PLN inquiry signature generated")
	
	// Verify signature manually for debugging
	expectedSignature := fmt.Sprintf("%x", md5.Sum([]byte(signatureInput)))
	if req.Sign != expectedSignature {
		c.logger.WithFields(logrus.Fields{
			"generated": req.Sign,
			"expected":  expectedSignature,
		}).Error("Signature mismatch detected!")
	} else {
		c.logger.Debug("Signature verification passed")
	}

	// Log request payload (without sensitive data)
	requestPayload := map[string]interface{}{
		"username":    req.Username,
		"customer_no": req.CustomerNo,
		"sign":        req.Sign,
	}
	c.logger.WithField("request_payload", requestPayload).Debug("PLN inquiry request payload")
	
	// Log full request for debugging
	fullRequest := map[string]interface{}{
		"username":    req.Username,
		"customer_no": req.CustomerNo,
		"sign":        req.Sign,
	}
	c.logger.WithField("full_request", fullRequest).Info("Full PLN inquiry request to Digiflazz")

	var resp models.PLNInquiryResponse
	if err := c.makeRequest("/inquiry-pln", req, &resp); err != nil {
		c.logger.WithError(err).Error("PLN inquiry request failed")
		return nil, err
	}

	// Log response details
	c.logger.WithFields(logrus.Fields{
		"customer_no": req.CustomerNo,
		"rc":          resp.Data.RC,
		"status":      resp.Data.Status,
		"message":     resp.Data.Message,
	}).Info("PLN inquiry response received from Digiflazz API")

	return &resp, nil
}

// generatePLNInquirySign generates MD5 signature for PLN inquiry
func (c *Client) generatePLNInquirySign(username, apiKey, customerNo string) string {
	data := fmt.Sprintf("%s%s%s", username, apiKey, customerNo)
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

// ValidateWebhook validates the webhook signature
func (c *Client) ValidateWebhook(webhook models.WebhookRequest) bool {
	expectedSign := c.generateSign(c.config.Username, c.config.APIKey, webhook.RefID)
	return webhook.Sign == expectedSign
}
