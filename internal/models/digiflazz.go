package models

import "time"

// DigiflazzRequest represents the base request structure for Digiflazz API
type DigiflazzRequest struct {
	Username string `json:"username"`
	APIKey   string `json:"api_key"`
	Sign     string `json:"sign"`
}

// DigiflazzResponse represents the base response structure from Digiflazz API
type DigiflazzResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
}

// BalanceRequest represents the request for checking balance
type BalanceRequest struct {
	DigiflazzRequest
}

// BalanceResponse represents the response for balance check
type BalanceResponse struct {
	Data struct {
		Deposit float64 `json:"deposit"`
	} `json:"data"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// PriceRequest represents the request for getting price list
type PriceRequest struct {
	DigiflazzRequest
	Type string `json:"type,omitempty"` // prabayar or pascabayar
}

// PriceResponse represents the response for price list
type PriceResponse struct {
	Data []Product `json:"data"`
	Message string `json:"message"`
	Status int     `json:"status"`
}

// Product represents a product in the price list
type Product struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	PriceType   string  `json:"price_type"`
	Status      string  `json:"status"`
	Description string  `json:"description"`
}

// TopupRequest represents the request for topup transaction
type TopupRequest struct {
	DigiflazzRequest
	BuyerSKU string `json:"buyer_sku"`
	CustomerNo string `json:"customer_no"`
	RefID     string `json:"ref_id"`
	Sign      string `json:"sign"`
}

// TopupResponse represents the response for topup transaction
type TopupResponse struct {
	Data struct {
		RefID       string  `json:"ref_id"`
		CustomerNo  string  `json:"customer_no"`
		BuyerSKU    string  `json:"buyer_sku"`
		Message     string  `json:"message"`
		RC          string  `json:"rc"`
		SN          string  `json:"sn"`
		BuyerLastSaldo float64 `json:"buyer_last_saldo"`
		BuyerSaldo  float64 `json:"buyer_saldo"`
		Price       float64 `json:"price"`
		Status      string  `json:"status"`
		Timestamp   string  `json:"timestamp"`
	} `json:"data"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// PayRequest represents the request for payment transaction
type PayRequest struct {
	DigiflazzRequest
	BuyerSKU   string `json:"buyer_sku"`
	CustomerNo string `json:"customer_no"`
	RefID      string `json:"ref_id"`
	Sign       string `json:"sign"`
}

// PayResponse represents the response for payment transaction
type PayResponse struct {
	Data struct {
		RefID       string  `json:"ref_id"`
		CustomerNo  string  `json:"customer_no"`
		BuyerSKU    string  `json:"buyer_sku"`
		Message     string  `json:"message"`
		RC          string  `json:"rc"`
		BuyerLastSaldo float64 `json:"buyer_last_saldo"`
		BuyerSaldo  float64 `json:"buyer_saldo"`
		Price       float64 `json:"price"`
		Status      string  `json:"status"`
		Timestamp   string  `json:"timestamp"`
	} `json:"data"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// StatusRequest represents the request for checking transaction status
type StatusRequest struct {
	DigiflazzRequest
	RefID string `json:"ref_id"`
}

// StatusResponse represents the response for transaction status
type StatusResponse struct {
	Data struct {
		RefID       string  `json:"ref_id"`
		CustomerNo  string  `json:"customer_no"`
		BuyerSKU    string  `json:"buyer_sku"`
		Message     string  `json:"message"`
		RC          string  `json:"rc"`
		SN          string  `json:"sn"`
		BuyerLastSaldo float64 `json:"buyer_last_saldo"`
		BuyerSaldo  float64 `json:"buyer_saldo"`
		Price       float64 `json:"price"`
		Status      string  `json:"status"`
		Timestamp   string  `json:"timestamp"`
	} `json:"data"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// WebhookRequest represents the webhook request from Digiflazz
type WebhookRequest struct {
	RefID       string  `json:"ref_id"`
	CustomerNo  string  `json:"customer_no"`
	BuyerSKU    string  `json:"buyer_sku"`
	Message     string  `json:"message"`
	RC          string  `json:"rc"`
	SN          string  `json:"sn"`
	BuyerLastSaldo float64 `json:"buyer_last_saldo"`
	BuyerSaldo  float64 `json:"buyer_saldo"`
	Price       float64 `json:"price"`
	Status      string  `json:"status"`
	Timestamp   string  `json:"timestamp"`
	Sign        string  `json:"sign"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Transaction represents a transaction record
type Transaction struct {
	ID          string    `json:"id"`
	RefID       string    `json:"ref_id"`
	CustomerNo  string    `json:"customer_no"`
	BuyerSKU    string    `json:"buyer_sku"`
	ProductName string    `json:"product_name"`
	Price       float64   `json:"price"`
	Status      string    `json:"status"`
	Message     string    `json:"message"`
	RC          string    `json:"rc"`
	SN          string    `json:"sn"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// APIError represents a Digiflazz API error
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}
