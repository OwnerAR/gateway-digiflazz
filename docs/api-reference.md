# API Reference

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
All API requests require proper Digiflazz credentials configured in the environment variables.

## Endpoints

### Health Check
```http
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "service": "digiflazz-gateway",
  "time": "2023-12-01T10:00:00Z"
}
```

### Balance

#### Get Balance
```http
GET /api/v1/balance
```

**Response:**
```json
{
  "success": true,
  "data": {
    "data": {
      "deposit": 1000000.00
    },
    "message": "success",
    "status": 1
  }
}
```

### Price List

#### Get Prices
```http
GET /api/v1/prices?type=prabayar
GET /api/v1/prices?type=pascabayar
```

**Query Parameters:**
- `type` (optional): Filter by type (`prabayar` or `pascabayar`)

**Response:**
```json
{
  "success": true,
  "data": {
    "data": [
      {
        "code": "pulsa10",
        "name": "Pulsa 10.000",
        "type": "prabayar",
        "category": "pulsa",
        "price": 10000,
        "price_type": "fixed",
        "status": "active",
        "description": "Pulsa 10.000"
      }
    ],
    "message": "success",
    "status": 1
  }
}
```

### Transactions

#### Topup
```http
POST /api/v1/transactions/topup
```

**Request Body:**
```json
{
  "ref_id": "TXN123456789",
  "customer_no": "08123456789",
  "buyer_sku": "pulsa10"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "data": {
      "ref_id": "TXN123456789",
      "customer_no": "08123456789",
      "buyer_sku": "pulsa10",
      "message": "Transaksi berhasil",
      "rc": "00",
      "sn": "1234567890",
      "buyer_last_saldo": 1000000,
      "buyer_saldo": 990000,
      "price": 10000,
      "status": "success",
      "timestamp": "2023-12-01 10:00:00"
    },
    "message": "success",
    "status": 1
  }
}
```

#### Pay (Pascabayar)
```http
POST /api/v1/transactions/pay
```

**Request Body:**
```json
{
  "ref_id": "TXN123456789",
  "customer_no": "12345678901",
  "buyer_sku": "pln20"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "data": {
      "ref_id": "TXN123456789",
      "customer_no": "12345678901",
      "buyer_sku": "pln20",
      "message": "Pembayaran berhasil",
      "rc": "00",
      "buyer_last_saldo": 1000000,
      "buyer_saldo": 980000,
      "price": 20000,
      "status": "success",
      "timestamp": "2023-12-01 10:00:00"
    },
    "message": "success",
    "status": 1
  }
}
```

#### Check Status
```http
GET /api/v1/transactions/{ref_id}/status
```

**Path Parameters:**
- `ref_id`: Transaction reference ID

**Response:**
```json
{
  "success": true,
  "data": {
    "data": {
      "ref_id": "TXN123456789",
      "customer_no": "08123456789",
      "buyer_sku": "pulsa10",
      "message": "Transaksi berhasil",
      "rc": "00",
      "sn": "1234567890",
      "buyer_last_saldo": 1000000,
      "buyer_saldo": 990000,
      "price": 10000,
      "status": "success",
      "timestamp": "2023-12-01 10:00:00"
    },
    "message": "success",
    "status": 1
  }
}
```

## Error Responses

All error responses follow this format:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error description",
    "details": "Additional error details"
  }
}
```

### Common Error Codes

- `INVALID_REQUEST`: Invalid request format
- `MISSING_REF_ID`: Missing reference ID
- `MISSING_CODE`: Missing product code
- `MISSING_CATEGORY`: Missing category
- `BALANCE_FAILED`: Failed to retrieve balance
- `PRICES_FAILED`: Failed to retrieve price list
- `TOPUP_FAILED`: Failed to process topup
- `PAYMENT_FAILED`: Failed to process payment
- `STATUS_CHECK_FAILED`: Failed to check transaction status
- `PRODUCT_NOT_FOUND`: Product not found
- `PRODUCTS_FAILED`: Failed to retrieve products
- `INVALID_WEBHOOK`: Invalid webhook format
- `WEBHOOK_FAILED`: Failed to process webhook

## Rate Limiting

The API implements rate limiting to prevent abuse. The default rate limit is 100 requests per minute per IP address.

## CORS

The API supports Cross-Origin Resource Sharing (CORS) for web applications. All origins are allowed by default.

## Webhooks

The gateway can receive webhooks from Digiflazz for transaction status updates. Webhook endpoints are automatically configured based on your Digiflazz account settings.

### Webhook Format

```json
{
  "ref_id": "TXN123456789",
  "customer_no": "08123456789",
  "buyer_sku": "pulsa10",
  "message": "Transaksi berhasil",
  "rc": "00",
  "sn": "1234567890",
  "buyer_last_saldo": 1000000,
  "buyer_saldo": 990000,
  "price": 10000,
  "status": "success",
  "timestamp": "2023-12-01 10:00:00",
  "sign": "abc123def456"
}
```
