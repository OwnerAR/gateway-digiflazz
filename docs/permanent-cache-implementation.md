# Permanent Cache Implementation for PLN Inquiry

## Overview

PLN inquiry data is now cached permanently without TTL (Time To Live) because PLN customer data is static and never changes.

## Why Permanent Cache?

### **PLN Data Characteristics:**
- **Customer Name**: Never changes
- **Meter Number**: Never changes  
- **Subscriber ID**: Never changes
- **Segment Power**: Never changes
- **Customer Number**: Never changes

### **Benefits:**
1. **Performance**: No need to call Digiflazz API repeatedly for same customer
2. **Cost Reduction**: Fewer API calls to Digiflazz
3. **Reliability**: Reduced dependency on external API
4. **Consistency**: Same customer data always returned

## Implementation Changes

### **1. Service Configuration**
```go
config: models.PLNInquiryConfig{
    CacheEnabled:   true,
    CacheTTL:       0, // No expiration for static PLN data
    CacheKeyPrefix: "pln_inquiry:",
}
```

### **2. Cache Storage Logic**
```go
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

// Use 0 TTL for permanent cache
return s.cache.Set(ctx, key, string(cacheData), 0)
```

### **3. Cache Retrieval Logic**
```go
// Check if cache is expired (only if ExpiresAt is set)
if !cache.ExpiresAt.IsZero() && time.Now().After(cache.ExpiresAt) {
    // Delete expired cache
    s.cache.Delete(ctx, key)
    return nil, fmt.Errorf("cache expired")
}
```

### **4. SQLite Cache Implementation**
```go
// Get - Handle permanent cache
query := `SELECT data FROM pln_inquiry_cache WHERE customer_no = ? AND (expires_at > ? OR expires_at = '0001-01-01 00:00:00')`

// Set - Handle TTL = 0
if ttl == 0 {
    // Permanent cache - use zero time
    expiresAt = time.Time{}
} else {
    expiresAt = now.Add(ttl)
}
```

## Cache Behavior

### **Cache Miss (First Request):**
1. Call Digiflazz API
2. Store response in permanent cache
3. Return response with proper ref_id and message

### **Cache Hit (Subsequent Requests):**
1. Retrieve data from permanent cache
2. Update ref_id to current request
3. Return cached data with new ref_id
4. No API call to Digiflazz

### **Cache Hit with Different Ref ID:**
1. Retrieve data from permanent cache
2. Update ref_id to new request value
3. Return same customer data with new ref_id
4. No API call to Digiflazz

## Database Schema

```sql
CREATE TABLE IF NOT EXISTS pln_inquiry_cache (
    customer_no TEXT PRIMARY KEY,
    data TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    expires_at DATETIME NOT NULL  -- '0001-01-01 00:00:00' for permanent cache
);

CREATE INDEX IF NOT EXISTS idx_expires_at ON pln_inquiry_cache(expires_at);
```

## Testing

### **1. Test Cache Behavior**
```bash
make test-cache-behavior
```

### **2. Manual Testing**
```bash
# Clear cache first
del data\cache.db

# Test first request (cache miss)
curl "http://localhost:8080/otomax/pln/inquiry?customer_no=543602392932&ref_id=12343"

# Test second request (cache hit - should be faster)
curl "http://localhost:8080/otomax/pln/inquiry?customer_no=543602392932&ref_id=12343"

# Test with different ref_id (cache hit with different ref_id)
curl "http://localhost:8080/otomax/pln/inquiry?customer_no=543602392932&ref_id=56789"
```

## Expected Results

### **Response 1 (Cache Miss):**
```json
{
  "data": {
    "data": {
      "message": "Transaksi Sukses",
      "status": "Sukses",
      "rc": "00",
      "ref_id": "12343",
      "customer_no": "543602392932",
      "meter_no": "32164069257",
      "subscriber_id": "543602392932",
      "name": "ZULFANTY YAHNEM",
      "segment_power": "R1   /2200"
    },
    "message": "Success",
    "status": 1
  },
  "message": "PLN inquiry completed successfully",
  "success": true,
  "timestamp": "2025-10-16T20:36:25Z"
}
```

### **Response 2 (Cache Hit):**
```json
{
  "data": {
    "data": {
      "message": "Transaksi Sukses",
      "status": "Sukses", 
      "rc": "00",
      "ref_id": "12343",  // Same ref_id
      "customer_no": "543602392932",
      "meter_no": "32164069257",
      "subscriber_id": "543602392932",
      "name": "ZULFANTY YAHNEM",
      "segment_power": "R1   /2200"
    },
    "message": "Success",
    "status": 1
  },
  "message": "PLN inquiry completed successfully (cached)",
  "success": true,
  "timestamp": "2025-10-16T20:36:25Z"
}
```

### **Response 3 (Cache Hit with Different Ref ID):**
```json
{
  "data": {
    "data": {
      "message": "Transaksi Sukses",
      "status": "Sukses",
      "rc": "00", 
      "ref_id": "56789",  // Different ref_id
      "customer_no": "543602392932",
      "meter_no": "32164069257",
      "subscriber_id": "543602392932",
      "name": "ZULFANTY YAHNEM",
      "segment_power": "R1   /2200"
    },
    "message": "Success",
    "status": 1
  },
  "message": "PLN inquiry completed successfully (cached)",
  "success": true,
  "timestamp": "2025-10-16T20:36:25Z"
}
```

## Cache Management

### **Clear Cache**
```bash
# Clear specific customer
curl -X DELETE "http://localhost:8080/pln/cache/543602392932"

# Clear all cache
curl -X DELETE "http://localhost:8080/pln/cache/all"
```

### **Cache Statistics**
```bash
# Get cache stats
curl "http://localhost:8080/pln/stats"
```

## Benefits

1. **Performance**: Subsequent requests are much faster
2. **Cost**: Reduced API calls to Digiflazz
3. **Reliability**: Less dependency on external API
4. **Consistency**: Same customer data always returned
5. **Scalability**: Better handling of high-volume requests

## Considerations

### **Cache Size Management**
- Monitor cache database size
- Implement cache cleanup if needed
- Consider cache compression for large datasets

### **Data Updates**
- If customer data changes (rare), manual cache clear required
- Consider implementing cache invalidation webhook from Digiflazz

### **Memory Usage**
- Permanent cache will grow over time
- Monitor SQLite database size
- Consider implementing cache size limits if needed

## Files Modified

- `internal/services/pln_inquiry.go`: Updated cache logic for permanent storage
- `pkg/cache/sqlite.go`: Updated SQLite implementation for permanent cache
- `scripts/test-cache-behavior.bat`: Updated test script for permanent cache behavior

The PLN inquiry cache is now optimized for static data with permanent storage and improved performance.
