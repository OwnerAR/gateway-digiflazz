# Pascabayar Transaction Flow

## Overview
Pascabayar (postpaid) transactions require a two-step process: **Check Bill** → **Pay Bill**. This ensures that the customer is aware of the exact amount before making the payment.

## Flow Diagram

```
Otomax → Gateway → Digiflazz
   ↓        ↓         ↓
1. Check Bill Request
   ↓        ↓         ↓
2. Bill Details Response
   ↓        ↓         ↓
3. Pay Bill Request
   ↓        ↓         ↓
4. Payment Confirmation
```

## API Endpoints

### 1. Check Bill (Pascabayar)

#### Standard API
```http
POST /api/v1/pascabayar/check
```

#### Otomax API (GET with query)
```http
GET /otomax/pascabayar/check
```

**Request Body (POST) / Query Parameters (GET):**
```json
{
  "ref_id": "TXN123456789",
  "customer_no": "12345678901",
  "buyer_sku": "pln20",
  "sign": "abc123def456"
}
```

**Response:**
```json
{
  "ref_id": "TXN123456789",
  "customer_no": "12345678901",
  "buyer_sku": "pln20",
  "amount": 50000,
  "admin_fee": 2500,
  "total": 52500,
  "status": "success",
  "message": "Bill check successful",
  "rc": "00",
  "bill_details": {
    "customer_name": "John Doe",
    "bill_period": "2023-12",
    "due_date": "2023-12-31",
    "bill_amount": 50000,
    "admin_fee": 2500,
    "total_amount": 52500
  },
  "timestamp": "2023-12-01T10:00:00Z",
  "sign": "def456ghi789"
}
```

### 2. Pay Bill (Pascabayar)

#### Standard API
```http
POST /api/v1/pascabayar/pay
```

#### Otomax API (GET with query)
```http
GET /otomax/pascabayar/pay
```

**Request Body (POST) / Query Parameters (GET):**
```json
{
  "ref_id": "TXN123456789",
  "customer_no": "12345678901",
  "buyer_sku": "pln20",
  "amount": 50000,
  "sign": "abc123def456"
}
```

**Response:**
```json
{
  "ref_id": "TXN123456789",
  "customer_no": "12345678901",
  "buyer_sku": "pln20",
  "amount": 50000,
  "admin_fee": 2500,
  "total": 52500,
  "status": "success",
  "message": "Bill payment successful",
  "rc": "00",
  "sn": "SN123456789",
  "bill_details": {
    "customer_name": "John Doe",
    "bill_period": "2023-12",
    "due_date": "2023-12-31",
    "bill_amount": 50000,
    "admin_fee": 2500,
    "total_amount": 52500
  },
  "timestamp": "2023-12-01T10:00:00Z",
  "sign": "def456ghi789"
}
```

## Transaction States

| State | Description |
|-------|-------------|
| `pending` | Transaction is being processed |
| `success` | Transaction completed successfully |
| `failed` | Transaction failed |

## Response Codes (RC)

| Code | Description |
|------|-------------|
| `00` | Success |
| `01` | Invalid request |
| `02` | Invalid signature |
| `03` | Bill not found |
| `04` | Insufficient balance |
| `05` | Bill already paid |
| `06` | Customer number invalid |
| `07` | Transaction timeout |
| `08` | System error |

## Bill Details Structure

```json
{
  "customer_name": "Customer Name",
  "bill_period": "YYYY-MM",
  "due_date": "YYYY-MM-DD",
  "bill_amount": 50000,
  "admin_fee": 2500,
  "total_amount": 52500
}
```

## Error Handling

### Common Error Responses

```json
{
  "code": "ERROR_CODE",
  "message": "Error description",
  "details": "Additional error details"
}
```

### Error Codes

- `INVALID_REQUEST`: Invalid request format
- `MISSING_PARAMETERS`: Missing required parameters
- `BILL_CHECK_FAILED`: Failed to check bill
- `BILL_PAYMENT_FAILED`: Failed to pay bill
- `TRANSACTION_NOT_FOUND`: Transaction not found

## Usage Examples

### 1. Check Bill (Standard API)
```bash
curl -X POST http://localhost:8080/api/v1/pascabayar/check \
  -H "Content-Type: application/json" \
  -d '{
    "ref_id": "TXN123456789",
    "customer_no": "12345678901",
    "buyer_sku": "pln20",
    "sign": "abc123def456"
  }'
```

### 2. Check Bill (Otomax API)
```bash
curl "http://localhost:8080/otomax/pascabayar/check?ref_id=TXN123456789&customer_no=12345678901&buyer_sku=pln20&sign=abc123def456"
```

### 3. Pay Bill (Standard API)
```bash
curl -X POST http://localhost:8080/api/v1/pascabayar/pay \
  -H "Content-Type: application/json" \
  -d '{
    "ref_id": "TXN123456789",
    "customer_no": "12345678901",
    "buyer_sku": "pln20",
    "amount": 50000,
    "sign": "abc123def456"
  }'
```

### 4. Pay Bill (Otomax API)
```bash
curl "http://localhost:8080/otomax/pascabayar/pay?ref_id=TXN123456789&customer_no=12345678901&buyer_sku=pln20&amount=50000&sign=abc123def456"
```

## Security Considerations

1. **Signature Validation**: All requests must include valid signatures
2. **Amount Validation**: Payment amount must match the checked bill amount
3. **Duplicate Prevention**: Prevent duplicate payments for the same bill
4. **Timeout Handling**: Implement proper timeout for bill checks
5. **Rate Limiting**: Limit the number of bill checks per customer

## Best Practices

1. **Always Check First**: Never process payment without checking the bill first
2. **Validate Amount**: Ensure payment amount matches the checked amount
3. **Handle Errors**: Implement proper error handling for failed checks
4. **Log Transactions**: Log all bill checks and payments for audit
5. **Monitor Performance**: Monitor response times for bill checks

## Integration Notes

- **Digiflazz API**: Uses `/pascabayar/check` and `/pascabayar/pay` endpoints
- **Signature Generation**: MD5 hash of specific parameters
- **Timeout**: Default 30 seconds for API calls
- **Retry Logic**: 3 attempts for failed requests
- **Webhook Support**: Callback notifications for status updates
