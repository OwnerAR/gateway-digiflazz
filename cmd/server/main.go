package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gateway-digiflazz/internal/config"
	"gateway-digiflazz/internal/handlers"
	"gateway-digiflazz/internal/middleware"
	"gateway-digiflazz/internal/services"
	"gateway-digiflazz/pkg/cache"
	"gateway-digiflazz/pkg/digiflazz"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// setDefaultEnvVars sets default environment variables if not already set
func setDefaultEnvVars() {
	defaults := map[string]string{
		"SERVER_HOST":         "0.0.0.0",
		"SERVER_PORT":         "8080",
		"SERVER_READ_TIMEOUT": "30s",
		"SERVER_WRITE_TIMEOUT": "30s",
		"SERVER_IDLE_TIMEOUT": "120s",
		"LOG_LEVEL":           "info",
		"LOG_FORMAT":          "text",
		"DIGIFLAZZ_USERNAME":  "your_username",
		"DIGIFLAZZ_API_KEY":   "your_api_key",
		"DIGIFLAZZ_BASE_URL":  "https://api.digiflazz.com",
		"OTOMAX_SECRET_KEY":   "default-secret-key",
		"OTOMAX_CALLBACK_URL": "http://localhost:8080/otomax/callback",
	}

	for key, value := range defaults {
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

// createEnvFileIfNotExists creates a .env file with default values if it doesn't exist
func createEnvFileIfNotExists() {
	envFile := ".env"
	
	// Check if .env file exists
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		log.Println("Creating .env file with default configuration...")
		
		// Create .env file with default values
		envContent := `# Digiflazz Gateway Configuration
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

# Database Configuration (if needed)
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
`

		// Write .env file
		if err := os.WriteFile(envFile, []byte(envContent), 0644); err != nil {
			log.Printf("Warning: Failed to create .env file: %v", err)
			log.Println("Application will use system environment variables and defaults")
		} else {
			log.Println("✅ .env file created successfully!")
			log.Println("📝 Please update the configuration values in .env file before using the application")
		}
	} else {
		log.Println("✅ .env file found, loading configuration...")
	}
}

func main() {
	// Check if .env file exists, if not create it
	createEnvFileIfNotExists()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Failed to load .env file:", err)
	}

	// Set default environment variables if not set
	setDefaultEnvVars()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	if cfg.Logging.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	// Initialize Digiflazz client
	digiflazzClient := digiflazz.NewClient(cfg.Digiflazz, logger)

	// Initialize SQLite cache
	sqliteCache, err := cache.NewSQLiteCache("cache.db")
	if err != nil {
		log.Fatalf("Failed to initialize SQLite cache: %v", err)
	}
	defer sqliteCache.Close()

	// Initialize services
	transactionService := services.NewTransactionService(digiflazzClient, logger)
	balanceService := services.NewBalanceService(digiflazzClient, logger)
	priceService := services.NewPriceService(digiflazzClient, logger)
	pascabayarService := services.NewPascabayarService(digiflazzClient, logger)
	plnInquiryService := services.NewPLNInquiryService(digiflazzClient, logger, sqliteCache)
	
	// Initialize Otomax service
	otomaxSecretKey := os.Getenv("OTOMAX_SECRET_KEY")
	if otomaxSecretKey == "" {
		otomaxSecretKey = "default-secret-key" // TODO: Use proper secret key management
	}
	otomaxService := services.NewOtomaxService(digiflazzClient, logger, otomaxSecretKey)

	// Initialize handlers
	transactionHandler := handlers.NewTransactionHandler(transactionService, logger)
	balanceHandler := handlers.NewBalanceHandler(balanceService, logger)
	priceHandler := handlers.NewPriceHandler(priceService, logger)
	pascabayarHandler := handlers.NewPascabayarHandler(pascabayarService, logger)
	plnInquiryHandler := handlers.NewPLNInquiryHandler(plnInquiryService, logger)
	otomaxHandler := handlers.NewOtomaxHandler(otomaxService, plnInquiryService, logger)

	// Setup router
	router := setupRouter(transactionHandler, balanceHandler, priceHandler, pascabayarHandler, plnInquiryHandler, otomaxHandler, logger)

	// Create server
	server := &http.Server{
		Addr:         cfg.Server.Host + ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in goroutine
	go func() {
		logger.Infof("Server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}

func setupRouter(
	transactionHandler *handlers.TransactionHandler,
	balanceHandler *handlers.BalanceHandler,
	priceHandler *handlers.PriceHandler,
	pascabayarHandler *handlers.PascabayarHandler,
	plnInquiryHandler *handlers.PLNInquiryHandler,
	otomaxHandler *handlers.OtomaxHandler,
	logger *logrus.Logger,
) *gin.Engine {
	// Set Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimit())
	router.Use(middleware.ResponseInterceptor(logger))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "digiflazz-gateway",
			"time":    time.Now().UTC(),
		})
	})

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Balance routes
		v1.GET("/balance", balanceHandler.GetBalance)

		// Price routes
		v1.GET("/prices", priceHandler.GetPrices)

		// Transaction routes
		transactions := v1.Group("/transactions")
		{
			transactions.POST("/topup", transactionHandler.Topup)
			transactions.POST("/pay", transactionHandler.Pay)
			transactions.GET("/:ref_id/status", transactionHandler.GetStatus)
		}

		// Pascabayar routes
		pascabayar := v1.Group("/pascabayar")
		{
			pascabayar.POST("/check", pascabayarHandler.CheckBill)
			pascabayar.POST("/pay", pascabayarHandler.PayBill)
			pascabayar.GET("/:ref_id", pascabayarHandler.GetTransaction)
		}

		// PLN Inquiry routes
		pln := v1.Group("/pln")
		{
			pln.POST("/inquiry", plnInquiryHandler.InquiryPLN)
			pln.GET("/stats", plnInquiryHandler.GetStats)
			pln.DELETE("/cache/:customer_no", plnInquiryHandler.ClearCache)
			pln.DELETE("/cache", plnInquiryHandler.ClearAllCache)
			pln.PUT("/cache/config", plnInquiryHandler.UpdateCacheConfig)
		}
	}

	// Otomax API routes (GET with query parameters)
	otomax := router.Group("/otomax")
	{
		// Transaction processing via GET with query parameters
		otomax.GET("/transaction", otomaxHandler.ProcessTransaction)
		
		// Status check via GET with query parameters
		otomax.GET("/status", otomaxHandler.CheckStatus)
		
		// Callback handling (POST for Digiflazz callbacks)
		otomax.POST("/callback", otomaxHandler.ProcessCallback)
		
		// Pascabayar endpoints for Otomax
		otomax.GET("/pascabayar/check", otomaxHandler.CheckPascabayarBill)
		otomax.GET("/pascabayar/pay", otomaxHandler.PayPascabayarBill)

		otomax.GET("/pln/inquiry", otomaxHandler.InquiryPLN)
		otomax.GET("/pln/stats", otomaxHandler.GetPLNStats)
		otomax.DELETE("/pln/cache/:customer_no", otomaxHandler.ClearPLNCache)
		otomax.DELETE("/pln/cache", otomaxHandler.ClearAllPLNCache)
		otomax.PUT("/pln/cache/config", otomaxHandler.UpdatePLNCacheConfig)
		
		// Additional Otomax endpoints
		otomax.GET("/history", otomaxHandler.GetTransactionHistory)
		otomax.GET("/products", otomaxHandler.GetProductList)
	}

	return router
}
