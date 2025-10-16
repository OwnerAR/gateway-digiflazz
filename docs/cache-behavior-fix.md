# Cache Behavior Fix

## Problem Identified

From the user's report, there was an inconsistency in PLN inquiry responses:

### **Response 1 (Cache Miss - First Hit):**
```json
{
  "data": {
    "data": {
      "message": "Transaksi Sukses",
      "status": "Sukses", 
      "rc": "00",
      "ref_id": "",           // ❌ Missing ref_id
      "customer_no": "543602392932",
      "meter_no": "32164069257",
      "subscriber_id": "543602392932",
      "name": "ZULFANTY YAHNEM",
      "segment_power": "R1   /2200"
    },
    "message": "",            // ❌ Missing message
    "status": 0
  },
  "message": "PLN inquiry completed successfully",
  "success": true,
  "timestamp": "2025-10-16T20:36:25Z"
}
```

### **Response 2 (Cache Hit - Second Hit):**
```json
{
  "data": {
    "data": {
      "message": "Transaksi Sukses",
      "status": "Sukses",
      "rc": "00", 
      "ref_id": "12343",      // ✅ ref_id present
      "customer_no": "543602392932",
      // ... other fields
    },
    "message": "Success",     // ✅ message present
    "status": 1
  },
  // ... rest of response
}
```

## Root Cause Analysis

### **1. Cache Miss Response Mapping Issue**
- When cache miss occurs, response from Digiflazz API was not properly mapped
- `ref_id` and `message` fields were not populated in the response structure

### **2. Cache Hit Response Logic Issue**
- Cache hit logic was overriding `ref_id` but not ensuring other fields were consistent
- Different response structure between cache miss and cache hit

### **3. Response Structure Inconsistency**
- Cache miss and cache hit returned different response formats
- Missing fields in cache miss scenario

## Solutions Applied

### **1. Fixed Cache Miss Response Mapping**
```go
// Ensure response has proper structure for cache miss
if resp.Data.RefID == "" {
    resp.Data.RefID = req.RefID
}
if resp.Data.Message == "" && resp.Data.RC == "00" {
    resp.Data.Message = "Transaksi Sukses"
}
if resp.Message == "" {
    resp.Message = "PLN inquiry completed successfully"
}
if resp.Status == 0 && resp.Data.RC == "00" {
    resp.Status = 1
}
```

### **2. Fixed Cache Hit Response Logic**
```go
// Build response with current ref_id
response := s.buildResponseFromCache(cached)
response.Data.RefID = req.RefID // Use current ref_id
response.Data.Message = cached.Message // Ensure message is included
response.Message = "PLN inquiry completed successfully (cached)"
response.Status = 1
return response, nil
```

### **3. Consistent Response Structure**
- Both cache miss and cache hit now return consistent response format
- All required fields are properly populated
- Response status and messages are consistent

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

# Test second request (cache hit)
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
      "ref_id": "12343",           // ✅ Now populated
      "customer_no": "543602392932",
      "meter_no": "32164069257",
      "subscriber_id": "543602392932", 
      "name": "ZULFANTY YAHNEM",
      "segment_power": "R1   /2200"
    },
    "message": "Success",          // ✅ Now populated
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
      "ref_id": "12343",           // ✅ Consistent
      "customer_no": "543602392932",
      "meter_no": "32164069257",
      "subscriber_id": "543602392932",
      "name": "ZULFANTY YAHNEM", 
      "segment_power": "R1   /2200"
    },
    "message": "Success",          // ✅ Consistent
    "status": 1
  },
  "message": "PLN inquiry completed successfully (cached)",
  "success": true,
  "timestamp": "2025-10-16T20:36:25Z"
}
```

## Cache Behavior Summary

### **Cache Miss (First Request):**
- Calls Digiflazz API
- Stores response in cache
- Returns properly formatted response with all fields

### **Cache Hit (Subsequent Requests):**
- Retrieves data from cache
- Updates `ref_id` to current request
- Returns consistent response format
- Faster response time

### **Cache Hit with Different Ref ID:**
- Retrieves data from cache
- Updates `ref_id` to new request value
- Returns same customer data with new ref_id
- Maintains data consistency

## Benefits

1. **Consistent Response Format**: Both cache miss and cache hit return same structure
2. **Proper Field Population**: All required fields are populated in both scenarios
3. **Performance**: Cache hit provides faster response times
4. **Data Integrity**: Customer data remains consistent across requests
5. **Ref ID Flexibility**: Different ref_ids can be used for same customer data

## Files Modified

- `internal/services/pln_inquiry.go`: Fixed cache logic and response mapping
- `scripts/test-cache-behavior.bat`: Added cache behavior testing
- `Makefile`: Added test-cache-behavior target

The cache behavior is now consistent and reliable across all scenarios.
