package models

import "time"

// OtomaxTransactionRequest represents the request from Otomax
type OtomaxTransactionRequest struct {
	RefID      string `form:"ref_id" json:"ref_id" binding:"required"`
	CustomerNo string `form:"customer_no" json:"customer_no" binding:"required"`
	BuyerSKU   string `form:"buyer_sku" json:"buyer_sku" binding:"required"`
	Amount     string `form:"amount" json:"amount"`
	Type       string `form:"type" json:"type"` // prabayar or pascabayar
	Timestamp  string `form:"timestamp" json:"timestamp"`
}

// OtomaxTransactionResponse represents the response to Otomax
type OtomaxTransactionResponse struct {
	RefID      string  `json:"ref_id"`
	CustomerNo string  `json:"customer_no"`
	BuyerSKU   string  `json:"buyer_sku"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"` // success, pending, failed
	Message    string  `json:"message"`
	RC         string  `json:"rc"`
	SN         string  `json:"sn,omitempty"`
	Timestamp  string  `json:"timestamp"`
	Sign       string  `json:"sign"`
}

// OtomaxStatusRequest represents status check request from Otomax
type OtomaxStatusRequest struct {
	RefID     string `form:"ref_id" json:"ref_id" binding:"required"`
	Timestamp string `form:"timestamp" json:"timestamp"`
}

// OtomaxStatusResponse represents status check response to Otomax
type OtomaxStatusResponse struct {
	RefID      string  `json:"ref_id"`
	CustomerNo string  `json:"customer_no"`
	BuyerSKU   string  `json:"buyer_sku"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	Message    string  `json:"message"`
	RC         string  `json:"rc"`
	SN         string  `json:"sn,omitempty"`
	Timestamp  string  `json:"timestamp"`
	Sign       string  `json:"sign"`
}

// OtomaxCallback represents the callback from Digiflazz for Otomax transactions
type OtomaxCallback struct {
	RefID      string  `json:"ref_id"`
	CustomerNo string  `json:"customer_no"`
	BuyerSKU   string  `json:"buyer_sku"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	Message    string  `json:"message"`
	RC         string  `json:"rc"`
	SN         string  `json:"sn,omitempty"`
	Timestamp  string  `json:"timestamp"`
	Sign       string  `json:"sign"`
}

// OtomaxTransaction represents an Otomax transaction record
type OtomaxTransaction struct {
	ID          string    `json:"id"`
	RefID       string    `json:"ref_id"`
	CustomerNo  string    `json:"customer_no"`
	BuyerSKU    string    `json:"buyer_sku"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Message     string    `json:"message"`
	RC          string    `json:"rc"`
	SN          string    `json:"sn"`
	DigiflazzRefID string `json:"digiflazz_ref_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OtomaxError represents an Otomax error response
type OtomaxError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
