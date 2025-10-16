# Automatic .env File Creation

The Digiflazz Gateway application now automatically creates a `.env` file with default configuration values when it starts for the first time.

## How It Works

1. **First Run**: When you start the application for the first time, it checks if a `.env` file exists
2. **Auto Creation**: If no `.env` file is found, it automatically creates one with default values
3. **Configuration**: The application then loads the configuration from the `.env` file

## Default Configuration

The automatically created `.env` file includes:

```bash
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
SERVER_IDLE_TIMEOUT=120s

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=text

# Digiflazz API Configuration
DIGIFLAZZ_USERNAME=your_username
DIGIFLAZZ_API_KEY=your_api_key
DIGIFLAZZ_BASE_URL=https://api.digiflazz.com

# Otomax Configuration
OTOMAX_SECRET_KEY=default-secret-key
OTOMAX_CALLBACK_URL=http://localhost:8080/otomax/callback

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=digiflazz_gateway
DB_USER=your_db_user
DB_PASSWORD=your_db_password

# Cache Configuration
CACHE_TYPE=sqlite
CACHE_SQLITE_PATH=cache.db
CACHE_TTL=24h

# Security Configuration
JWT_SECRET=your-jwt-secret-key
API_KEY=your-api-key

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m

# Health Check
HEALTH_CHECK_INTERVAL=30s
HEALTH_CHECK_TIMEOUT=5s
```

## Usage

### First Time Setup

1. **Run the application**:
   ```bash
   go run ./cmd/server
   ```

2. **Check the output**:
   ```
   Creating .env file with default configuration...
   ✅ .env file created successfully!
   📝 Please update the configuration values in .env file before using the application
   ```

3. **Update configuration**:
   - Edit the `.env` file with your actual values
   - Replace `your_username`, `your_api_key`, etc. with real values

### Subsequent Runs

After the first run, the application will:
```
✅ .env file found, loading configuration...
```

## Customization

You can customize the default values by modifying the `createEnvFileIfNotExists()` function in `cmd/server/main.go`.

## Benefits

- **Zero Configuration**: No need to manually create `.env` file
- **Default Values**: Sensible defaults for all configuration options
- **Easy Setup**: Just run the application and it handles the rest
- **Production Ready**: Clear separation between development and production configs

## Notes

- The `.env` file is created in the same directory as the application
- If the application can't create the `.env` file, it will use system environment variables and defaults
- The `.env` file is not overwritten if it already exists
- Always update the configuration values before using in production


