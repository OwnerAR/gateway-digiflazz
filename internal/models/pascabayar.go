package models

import "time"

// PascabayarCheckRequest represents the request for checking Pascabayar bill
type PascabayarCheckRequest struct {
	DigiflazzRequest
	BuyerSKU   string `json:"buyer_sku" binding:"required"`
	CustomerNo string `json:"customer_no" binding:"required"`
	RefID      string `json:"ref_id" binding:"required"`
	Sign       string `json:"sign" binding:"required"`
}

// PascabayarCheckResponse represents the response for Pascabayar bill check
type PascabayarCheckResponse struct {
	Data struct {
		RefID      string  `json:"ref_id"`
		CustomerNo string  `json:"customer_no"`
		BuyerSKU   string  `json:"buyer_sku"`
		Message    string  `json:"message"`
		RC         string  `json:"rc"`
		Amount     float64 `json:"amount"`
		AdminFee   float64 `json:"admin_fee"`
		Total      float64 `json:"total"`
		Status     string  `json:"status"`
		Timestamp  string  `json:"timestamp"`
		BillDetails struct {
			CustomerName string `json:"customer_name"`
			BillPeriod   string `json:"bill_period"`
			DueDate      string `json:"due_date"`
			BillAmount   float64 `json:"bill_amount"`
			AdminFee     float64 `json:"admin_fee"`
			TotalAmount  float64 `json:"total_amount"`
		} `json:"bill_details"`
	} `json:"data"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// PascabayarPayRequest represents the request for paying Pascabayar bill
type PascabayarPayRequest struct {
	DigiflazzRequest
	BuyerSKU   string  `json:"buyer_sku" binding:"required"`
	CustomerNo string  `json:"customer_no" binding:"required"`
	RefID      string  `json:"ref_id" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	Sign       string  `json:"sign" binding:"required"`
}

// PascabayarPayResponse represents the response for Pascabayar bill payment
type PascabayarPayResponse struct {
	Data struct {
		RefID      string  `json:"ref_id"`
		CustomerNo string  `json:"customer_no"`
		BuyerSKU   string  `json:"buyer_sku"`
		Message    string  `json:"message"`
		RC         string  `json:"rc"`
		Amount     float64 `json:"amount"`
		AdminFee   float64 `json:"admin_fee"`
		Total      float64 `json:"total"`
		Status     string  `json:"status"`
		Timestamp  string  `json:"timestamp"`
		SN         string  `json:"sn"`
		BillDetails struct {
			CustomerName string `json:"customer_name"`
			BillPeriod   string `json:"bill_period"`
			DueDate      string `json:"due_date"`
			BillAmount   float64 `json:"bill_amount"`
			AdminFee     float64 `json:"admin_fee"`
			TotalAmount  float64 `json:"total_amount"`
		} `json:"bill_details"`
	} `json:"data"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// OtomaxPascabayarCheckRequest represents Otomax request for Pascabayar check
type OtomaxPascabayarCheckRequest struct {
	RefID      string `form:"ref_id" json:"ref_id" binding:"required"`
	CustomerNo string `form:"customer_no" json:"customer_no" binding:"required"`
	BuyerSKU   string `form:"buyer_sku" json:"buyer_sku" binding:"required"`
	Timestamp  string `form:"timestamp" json:"timestamp"`
}

// OtomaxPascabayarCheckResponse represents Otomax response for Pascabayar check
type OtomaxPascabayarCheckResponse struct {
	RefID        string  `json:"ref_id"`
	CustomerNo   string  `json:"customer_no"`
	BuyerSKU     string  `json:"buyer_sku"`
	Amount       float64 `json:"amount"`
	AdminFee     float64 `json:"admin_fee"`
	Total        float64 `json:"total"`
	Status       string  `json:"status"`
	Message      string  `json:"message"`
	RC           string  `json:"rc"`
	BillDetails  BillDetails `json:"bill_details"`
	Timestamp    string  `json:"timestamp"`
	Sign         string  `json:"sign"`
}

// OtomaxPascabayarPayRequest represents Otomax request for Pascabayar payment
type OtomaxPascabayarPayRequest struct {
	RefID      string  `form:"ref_id" json:"ref_id" binding:"required"`
	CustomerNo string  `form:"customer_no" json:"customer_no" binding:"required"`
	BuyerSKU   string  `form:"buyer_sku" json:"buyer_sku" binding:"required"`
	Amount     float64 `form:"amount" json:"amount" binding:"required"`
	Timestamp  string  `form:"timestamp" json:"timestamp"`
}

// OtomaxPascabayarPayResponse represents Otomax response for Pascabayar payment
type OtomaxPascabayarPayResponse struct {
	RefID        string  `json:"ref_id"`
	CustomerNo   string  `json:"customer_no"`
	BuyerSKU     string  `json:"buyer_sku"`
	Amount       float64 `json:"amount"`
	AdminFee     float64 `json:"admin_fee"`
	Total        float64 `json:"total"`
	Status       string  `json:"status"`
	Message      string  `json:"message"`
	RC           string  `json:"rc"`
	SN           string  `json:"sn"`
	BillDetails  BillDetails `json:"bill_details"`
	Timestamp    string  `json:"timestamp"`
	Sign         string  `json:"sign"`
}

// BillDetails represents bill details for Pascabayar transactions
type BillDetails struct {
	CustomerName string  `json:"customer_name"`
	BillPeriod   string  `json:"bill_period"`
	DueDate      string  `json:"due_date"`
	BillAmount   float64 `json:"bill_amount"`
	AdminFee     float64 `json:"admin_fee"`
	TotalAmount  float64 `json:"total_amount"`
}

// PascabayarTransaction represents a Pascabayar transaction record
type PascabayarTransaction struct {
	ID            string       `json:"id"`
	RefID         string       `json:"ref_id"`
	CustomerNo    string       `json:"customer_no"`
	BuyerSKU      string       `json:"buyer_sku"`
	Amount        float64      `json:"amount"`
	AdminFee      float64      `json:"admin_fee"`
	Total         float64      `json:"total"`
	Status        string       `json:"status"`
	Message       string       `json:"message"`
	RC            string       `json:"rc"`
	SN            string       `json:"sn"`
	BillDetails   BillDetails  `json:"bill_details"`
	DigiflazzRefID string     `json:"digiflazz_ref_id"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}
