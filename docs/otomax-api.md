# Otomax API Integration

## Overview
This document describes the API integration between Otomax and the Digiflazz Gateway for transaction processing using HTTP GET requests with query parameters.

## Base URL
```
http://localhost:8080/otomax
```

## Authentication
Requests from Otomax do not require signature validation. The gateway handles all signature generation and validation for Digiflazz API calls internally.

## Endpoints

### 1. Process Transaction
```http
GET /otomax/transaction
```

**Query Parameters:**
- `ref_id` (required): Unique transaction reference ID
- `customer_no` (required): Customer phone number or account number
- `buyer_sku` (required): Product SKU code
- `amount` (optional): Transaction amount
- `type` (optional): Transaction type (`prabayar` or `pascabayar`)
- `timestamp` (optional): Request timestamp

**Example Request:**
```
GET /otomax/transaction?ref_id=TXN123456789&customer_no=08123456789&buyer_sku=pulsa10&type=prabayar&timestamp=2023-12-01T10:00:00Z
```

**Response:**
```json
{
  "ref_id": "TXN123456789",
  "customer_no": "08123456789",
  "buyer_sku": "pulsa10",
  "amount": 10000,
  "status": "success",
  "message": "Transaksi berhasil",
  "rc": "00",
  "sn": "1234567890",
  "timestamp": "2023-12-01T10:00:00Z",
  "sign": "def456ghi789"
}
```

### 2. Check Transaction Status
```http
GET /otomax/status
```

**Query Parameters:**
- `ref_id` (required): Transaction reference ID
- `timestamp` (optional): Request timestamp

**Example Request:**
```
GET /otomax/status?ref_id=TXN123456789&timestamp=2023-12-01T10:00:00Z
```

**Response:**
```json
{
  "ref_id": "TXN123456789",
  "customer_no": "08123456789",
  "buyer_sku": "pulsa10",
  "amount": 10000,
  "status": "success",
  "message": "Transaction completed",
  "rc": "00",
  "sn": "1234567890",
  "timestamp": "2023-12-01T10:00:00Z",
  "sign": "def456ghi789"
}
```

### 3. Process Callback (Webhook)
```http
POST /otomax/callback
```

**Request Body:**
```json
{
  "ref_id": "TXN123456789",
  "customer_no": "08123456789",
  "buyer_sku": "pulsa10",
  "amount": 10000,
  "status": "success",
  "message": "Transaksi berhasil",
  "rc": "00",
  "sn": "1234567890",
  "timestamp": "2023-12-01T10:00:00Z",
  "sign": "def456ghi789"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Callback processed successfully"
}
```

### 4. Transaction History
```http
GET /otomax/history
```

**Response:**
```json
{
  "success": true,
  "message": "Transaction history endpoint - to be implemented",
  "data": []
}
```

### 5. Product List
```http
GET /otomax/products
```

**Response:**
```json
{
  "success": true,
  "message": "Product list endpoint - to be implemented",
  "data": []
}
```

## Signature Generation

**Note:** Signature generation is handled internally by the gateway for Digiflazz API calls. Otomax requests do not require signature validation.

### For Digiflazz API Calls (Internal)
The gateway automatically generates signatures for all Digiflazz API calls using the configured API key and username.

### For Response Signatures (Optional)
Response signatures are generated for consistency but are not required for basic functionality.

## Status Codes

| Status | Description |
|--------|-------------|
| `success` | Transaction completed successfully |
| `pending` | Transaction is being processed |
| `failed` | Transaction failed |

## Response Codes (RC)

| Code | Description |
|------|-------------|
| `00` | Success |
| `01` | Invalid request |
| `02` | Invalid signature |
| `03` | Transaction not found |
| `04` | Insufficient balance |
| `05` | Product not available |
| `06` | Customer number invalid |
| `07` | Transaction timeout |
| `08` | System error |

## Error Responses

All error responses follow this format:

```json
{
  "code": "ERROR_CODE",
  "message": "Error description",
  "details": "Additional error details"
}
```

### Common Error Codes

- `INVALID_REQUEST`: Invalid request parameters
- `MISSING_PARAMETERS`: Missing required parameters
- `TRANSACTION_FAILED`: Failed to process transaction
- `STATUS_CHECK_FAILED`: Failed to check transaction status
- `INVALID_CALLBACK`: Invalid callback format
- `CALLBACK_FAILED`: Failed to process callback

## Example Integration

### 1. Process Prabayar Transaction
```bash
curl "http://localhost:8080/otomax/transaction?ref_id=TXN001&customer_no=08123456789&buyer_sku=pulsa10&type=prabayar&timestamp=2023-12-01T10:00:00Z"
```

### 2. Process Pascabayar Transaction
```bash
curl "http://localhost:8080/otomax/transaction?ref_id=TXN002&customer_no=12345678901&buyer_sku=pln20&type=pascabayar&timestamp=2023-12-01T10:00:00Z"
```

### 3. Check Transaction Status
```bash
curl "http://localhost:8080/otomax/status?ref_id=TXN001&timestamp=2023-12-01T10:00:00Z"
```

## Security Considerations

1. **Internal Signature Handling**: Gateway handles all Digiflazz API signatures internally
2. **Timestamp Validation**: Consider implementing timestamp validation to prevent replay attacks
3. **Rate Limiting**: Implement rate limiting to prevent abuse
4. **IP Whitelisting**: Consider whitelisting Otomax IP addresses
5. **HTTPS**: Use HTTPS for all communications in production

## Configuration

Set the following environment variables:

```bash
OTOMAX_SECRET_KEY=your_secret_key_here
OTOMAX_CALLBACK_URL=https://your-domain.com/otomax/callback
```

## Testing

Use the following test cases:

1. **Valid Transaction**: Test with valid parameters and signature
2. **Invalid Signature**: Test with invalid signature
3. **Missing Parameters**: Test with missing required parameters
4. **Invalid Product**: Test with non-existent product SKU
5. **Network Error**: Test with network connectivity issues
