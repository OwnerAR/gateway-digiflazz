# PLN Inquiry API with SQLite Caching

## Overview
PLN Inquiry API allows you to check PLN customer information with intelligent caching using SQLite database. The cache strategy is based on customer_id from the request, providing fast response times and reduced API calls to Digiflazz.

## Base URL
```
http://localhost:8080/api/v1/pln
```

## Cache Strategy

### SQLite Database Schema
```sql
CREATE TABLE pln_inquiry_cache (
    customer_no TEXT PRIMARY KEY,
    data TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    expires_at DATETIME NOT NULL
);
```

### Cache Configuration
- **Default TTL**: 24 hours
- **Cache Key**: `pln_inquiry:{customer_no}`
- **Auto-cleanup**: Expired entries are automatically removed
- **Storage**: SQLite database (`cache.db`)

## API Endpoints

### 1. PLN Inquiry
```http
POST /api/v1/pln/inquiry
```

**Request Body:**
```json
{
  "customer_no": "1234554321",
  "sign": "740b00a1b8784e028cc8078edf66d12b"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "data": {
      "message": "Transaksi Sukses",
      "status": "Sukses",
      "rc": "00",
      "customer_no": "1234554321",
      "meter_no": "1234554321",
      "subscriber_id": "523300817840",
      "name": "DAVID",
      "segment_power": "R1 /000001300"
    },
    "message": "Success",
    "status": 1
  }
}
```

### 2. Cache Statistics
```http
GET /api/v1/pln/stats
```

**Response:**
```json
{
  "success": true,
  "data": {
    "total_requests": 150,
    "cache_hits": 120,
    "cache_misses": 30,
    "api_requests": 30,
    "error_count": 2,
    "average_response_time": "150ms"
  }
}
```

### 3. Clear Cache for Specific Customer
```http
DELETE /api/v1/pln/cache/{customer_no}
```

**Response:**
```json
{
  "success": true,
  "message": "Cache cleared successfully"
}
```

### 4. Clear All Cache
```http
DELETE /api/v1/pln/cache
```

**Response:**
```json
{
  "success": true,
  "message": "All cache cleared successfully"
}
```

### 5. Update Cache Configuration
```http
PUT /api/v1/pln/cache/config
```

**Request Body:**
```json
{
  "cache_enabled": true,
  "cache_ttl": "24h",
  "cache_key_prefix": "pln_inquiry:"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Cache configuration updated successfully",
  "data": {
    "cache_enabled": true,
    "cache_ttl": "24h",
    "cache_key_prefix": "pln_inquiry:"
  }
}
```

## Cache Behavior

### Cache Hit Flow
```
1. Request PLN inquiry for customer_no
2. Check SQLite cache for existing data
3. If found and not expired → Return cached data
4. Log cache hit statistics
```

### Cache Miss Flow
```
1. Request PLN inquiry for customer_no
2. Check SQLite cache for existing data
3. If not found or expired → Call Digiflazz API
4. Store response in SQLite cache
5. Return API response
6. Log cache miss statistics
```

## Signature Generation

The signature is generated using MD5 hash:
```
sign = MD5(username + apiKey + customer_no)
```

## Response Codes

| Code | Description |
|------|-------------|
| `00` | Success |
| `01` | Invalid request |
| `02` | Invalid signature |
| `03` | Customer not found |
| `04` | System error |

## Error Responses

```json
{
  "code": "ERROR_CODE",
  "message": "Error description",
  "details": "Additional error details"
}
```

### Common Error Codes

- `INVALID_REQUEST`: Invalid request format
- `MISSING_CUSTOMER_NO`: Missing customer_no parameter
- `INQUIRY_FAILED`: Failed to perform PLN inquiry
- `CACHE_CLEAR_FAILED`: Failed to clear cache

## Usage Examples

### 1. Basic PLN Inquiry
```bash
curl -X POST http://localhost:8080/api/v1/pln/inquiry \
  -H "Content-Type: application/json" \
  -d '{
    "customer_no": "1234554321",
    "sign": "740b00a1b8784e028cc8078edf66d12b"
  }'
```

### 2. Check Cache Statistics
```bash
curl http://localhost:8080/api/v1/pln/stats
```

### 3. Clear Specific Customer Cache
```bash
curl -X DELETE http://localhost:8080/api/v1/pln/cache/1234554321
```

### 4. Clear All Cache
```bash
curl -X DELETE http://localhost:8080/api/v1/pln/cache
```

### 5. Update Cache Configuration
```bash
curl -X PUT http://localhost:8080/api/v1/pln/cache/config \
  -H "Content-Type: application/json" \
  -d '{
    "cache_enabled": true,
    "cache_ttl": "12h",
    "cache_key_prefix": "pln_inquiry:"
  }'
```

## Performance Benefits

### Cache Hit Benefits
- **Response Time**: ~5ms (vs ~200ms API call)
- **API Calls**: Reduced by 80-90%
- **Cost**: Lower Digiflazz API usage
- **Reliability**: Works even if Digiflazz is temporarily down

### Cache Statistics
- **Cache Hit Rate**: Typically 80-90%
- **Average Response Time**: 5ms (cached) vs 200ms (API)
- **Storage**: Minimal SQLite database size
- **Cleanup**: Automatic expired entry removal

## Database Management

### SQLite Database Location
- **File**: `cache.db` (in application directory)
- **Size**: Typically < 10MB for 10,000 cached entries
- **Backup**: Regular SQLite file backup recommended

### Cache Cleanup
```sql
-- Manual cleanup of expired entries
DELETE FROM pln_inquiry_cache WHERE expires_at <= datetime('now');

-- Check cache statistics
SELECT 
    COUNT(*) as total_entries,
    COUNT(CASE WHEN expires_at > datetime('now') THEN 1 END) as active_entries,
    COUNT(CASE WHEN expires_at <= datetime('now') THEN 1 END) as expired_entries
FROM pln_inquiry_cache;
```

## Monitoring and Maintenance

### Key Metrics to Monitor
1. **Cache Hit Rate**: Should be > 80%
2. **Average Response Time**: Should be < 10ms for cached requests
3. **Database Size**: Monitor SQLite file size
4. **Error Rate**: Should be < 1%

### Maintenance Tasks
1. **Regular Cleanup**: Remove expired entries
2. **Database Backup**: Backup SQLite file
3. **Performance Monitoring**: Track cache statistics
4. **Configuration Updates**: Adjust TTL based on usage patterns

## Security Considerations

1. **Database Security**: Secure SQLite file permissions
2. **Signature Validation**: All requests must be properly signed
3. **Rate Limiting**: Implement rate limiting for API endpoints
4. **Data Privacy**: Cache contains customer information - ensure compliance
5. **Access Control**: Restrict cache management endpoints

## Troubleshooting

### Common Issues

1. **Cache Not Working**
   - Check SQLite database permissions
   - Verify cache configuration
   - Check database file exists

2. **High Memory Usage**
   - Monitor SQLite database size
   - Implement cache size limits
   - Regular cleanup of expired entries

3. **Slow Performance**
   - Check database indexes
   - Monitor cache hit rate
   - Optimize SQLite configuration

### Debug Commands
```bash
# Check cache statistics
curl http://localhost:8080/api/v1/pln/stats

# Clear all cache
curl -X DELETE http://localhost:8080/api/v1/pln/cache

# Test PLN inquiry
curl -X POST http://localhost:8080/api/v1/pln/inquiry \
  -H "Content-Type: application/json" \
  -d '{"customer_no": "1234554321", "sign": "test"}'
```
