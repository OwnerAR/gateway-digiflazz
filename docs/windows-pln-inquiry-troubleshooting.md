# Windows PLN Inquiry Troubleshooting Guide

## Masalah: PLN Inquiry Response Kosong di Windows

### Gejala
- PLN inquiry berjalan di macOS/Linux tapi response kosong di Windows
- Log menunjukkan `retry_attempts=0 timeout=0s user_agent="<nil>"`
- Digiflazz API call berhasil tapi response kosong

### Root Cause
Masalah utama adalah **environment variables tidak ter-load dengan benar di Windows**, sehingga konfigurasi Digiflazz client tidak optimal.

### Solusi

#### 1. Set Environment Variables Secara Manual
```batch
REM Set proper environment variables sebelum menjalankan aplikasi
set DIGIFLAZZ_TIMEOUT=30s
set DIGIFLAZZ_RETRY_ATTEMPTS=3
set DIGIFLAZZ_BASE_URL=https://api.digiflazz.com
set DIGIFLAZZ_IP_WHITELIST=52.74.250.133
```

#### 2. Gunakan Script Testing yang Sudah Disiapkan
```bash
# Debug konfigurasi
make debug-config-windows

# Test PLN inquiry dengan konfigurasi yang benar
make test-pln-inquiry-windows
```

#### 3. Verifikasi .env File
Pastikan file `.env` ada dan berisi konfigurasi yang benar:
```env
DIGIFLAZZ_USERNAME=your_username
DIGIFLAZZ_API_KEY=your_api_key
DIGIFLAZZ_BASE_URL=https://api.digiflazz.com
DIGIFLAZZ_TIMEOUT=30s
DIGIFLAZZ_RETRY_ATTEMPTS=3
DIGIFLAZZ_IP_WHITELIST=52.74.250.133
```

### Cache Strategy PLN Inquiry

#### Alur PLN Inquiry dari Otomax:
1. **Cache Check First**: `getFromCache(customerNo)` dipanggil terlebih dahulu
2. **Cache Hit**: Jika data ada di cache, langsung return response dengan `ref_id` dari request Otomax
3. **Cache Miss**: Jika tidak ada di cache, hit ke Digiflazz API
4. **Cache Store**: Response dari API disimpan ke cache jika sukses (`RC == "00"`)
5. **Response Mapping**: Response di-mapping dengan `ref_id` dari request Otomax

#### Permanent Cache Implementation:
- PLN inquiry data bersifat static (tidak berubah)
- Cache TTL diset ke 0 (permanent cache)
- Cache key: `pln_inquiry:{customer_no}`
- Cache storage: SQLite dengan `expires_at = '0001-01-01 00:00:00'`

### Debugging Steps

#### 1. Check Configuration
```batch
scripts/debug-config-windows.bat
```

#### 2. Test PLN Inquiry
```batch
scripts/test-pln-inquiry-windows.bat
```

#### 3. Check Cache Behavior
```batch
scripts/test-cache-behavior.bat
```

#### 4. PowerShell Debug
```powershell
powershell -ExecutionPolicy Bypass -File scripts/debug-windows.ps1
```

### Expected Log Output (Normal)
```
time="2025-10-17T03:19:36+07:00" level=info msg="Digiflazz client configuration" 
  api_key_len=36 base_url="https://api.digiflazz.com/v1" 
  platform=windows/amd64 retry_attempts=3 timeout=30s 
  user_agent="Digiflazz-Gateway/1.0 (windows/amd64)" username=zokafio0jV7W
```

### Troubleshooting Checklist

- [ ] Environment variables diset dengan benar
- [ ] .env file ada dan readable
- [ ] Digiflazz credentials valid
- [ ] IP address ter-whitelist di Digiflazz
- [ ] Network connectivity ke api.digiflazz.com
- [ ] Cache database path accessible
- [ ] Application running dengan proper permissions

### Common Solutions

#### 1. Restart Application dengan Environment Variables
```batch
# Set environment variables
set DIGIFLAZZ_TIMEOUT=30s
set DIGIFLAZZ_RETRY_ATTEMPTS=3

# Start application
gateway-digiflazz-windows-amd64.exe
```

#### 2. Clear Cache Database
```batch
# Delete cache database to force fresh API calls
del data\cache.db
```

#### 3. Test dengan Customer Number yang Valid
```batch
curl -X GET "http://localhost:8081/otomax/pln/inquiry?customer_no=543602392932&ref_id=12343"
```

### Performance Monitoring

#### Cache Statistics
- `CacheHits`: Jumlah request yang dilayani dari cache
- `CacheMisses`: Jumlah request yang harus hit ke Digiflazz API
- `APIRequests`: Total request ke Digiflazz API
- `AverageResponseTime`: Rata-rata waktu response

#### Expected Behavior
- **First Hit**: Cache miss → API call → Store to cache
- **Subsequent Hits**: Cache hit → Serve from cache
- **Response Time**: Cache hit < 10ms, API call ~500-2000ms

### References
- [Digiflazz API Documentation](https://developer.digiflazz.com/api/)
- [Windows Environment Variables](https://docs.microsoft.com/en-us/windows/win32/procthread/environment-variables)
- [SQLite Cache Implementation](./cache-behavior-fix.md)
- [Permanent Cache Implementation](./permanent-cache-implementation.md)
