package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the application
type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Digiflazz  DigiflazzConfig  `yaml:"digiflazz"`
	Database   DatabaseConfig   `yaml:"database"`
	Redis      RedisConfig      `yaml:"redis"`
	Logging    LoggingConfig    `yaml:"logging"`
	Security   SecurityConfig   `yaml:"security"`
	Monitoring MonitoringConfig `yaml:"monitoring"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// DigiflazzConfig holds Digiflazz API configuration
type DigiflazzConfig struct {
	BaseURL      string        `yaml:"base_url"`
	Username     string        `yaml:"username"`
	APIKey       string        `yaml:"api_key"`
	IPWhitelist  string        `yaml:"ip_whitelist"`
	Timeout      time.Duration `yaml:"timeout"`
	RetryAttempts int          `yaml:"retry_attempts"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host               string        `yaml:"host"`
	Port               int           `yaml:"port"`
	Name               string        `yaml:"name"`
	User               string        `yaml:"user"`
	Password           string        `yaml:"password"`
	SSLMode            string        `yaml:"ssl_mode"`
	MaxConnections     int           `yaml:"max_connections"`
	MaxIdleConnections int           `yaml:"max_idle_connections"`
	ConnectionMaxLifetime time.Duration `yaml:"connection_max_lifetime"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host             string `yaml:"host"`
	Port             int    `yaml:"port"`
	Password         string `yaml:"password"`
	DB               int    `yaml:"db"`
	PoolSize         int    `yaml:"pool_size"`
	MinIdleConnections int  `yaml:"min_idle_connections"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Output     string `yaml:"output"`
	FilePath   string `yaml:"file_path"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
}

// SecurityConfig holds security configuration
type SecurityConfig struct {
	JWTSecret   string   `yaml:"jwt_secret"`
	APIRateLimit int     `yaml:"api_rate_limit"`
	CORSOrigins []string `yaml:"cors_origins"`
	CORSMethods []string `yaml:"cors_methods"`
	CORSHeaders []string `yaml:"cors_headers"`
}

// MonitoringConfig holds monitoring configuration
type MonitoringConfig struct {
	EnableMetrics        bool          `yaml:"enable_metrics"`
	MetricsPort          int           `yaml:"metrics_port"`
	HealthCheckInterval  time.Duration `yaml:"health_check_interval"`
}

// Load loads configuration from environment variables and config file
func Load() (*Config, error) {
	cfg := &Config{}

	// Load from YAML file if exists
	if err := loadFromYAML(cfg); err != nil {
		return nil, fmt.Errorf("failed to load YAML config: %w", err)
	}

	// Override with environment variables
	loadFromEnv(cfg)

	return cfg, nil
}

// loadFromYAML loads configuration from YAML file
func loadFromYAML(cfg *Config) error {
	configFile := "configs/config.yaml"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Use default configuration if file doesn't exist
		return nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, cfg)
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv(cfg *Config) {
	// Server configuration
	if host := os.Getenv("SERVER_HOST"); host != "" {
		cfg.Server.Host = host
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.Server.Port = port
	}

	// Digiflazz configuration
	if baseURL := os.Getenv("DIGIFLAZZ_BASE_URL"); baseURL != "" {
		cfg.Digiflazz.BaseURL = baseURL
	}
	if username := os.Getenv("DIGIFLAZZ_USERNAME"); username != "" {
		cfg.Digiflazz.Username = username
	}
	if apiKey := os.Getenv("DIGIFLAZZ_API_KEY"); apiKey != "" {
		cfg.Digiflazz.APIKey = apiKey
	}
	if ipWhitelist := os.Getenv("DIGIFLAZZ_IP_WHITELIST"); ipWhitelist != "" {
		cfg.Digiflazz.IPWhitelist = ipWhitelist
	}
	
	// Timeout configuration
	if timeoutStr := os.Getenv("DIGIFLAZZ_TIMEOUT"); timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil {
			cfg.Digiflazz.Timeout = timeout
		}
	}
	
	// Retry attempts configuration
	if retryStr := os.Getenv("DIGIFLAZZ_RETRY_ATTEMPTS"); retryStr != "" {
		if retry, err := strconv.Atoi(retryStr); err == nil {
			cfg.Digiflazz.RetryAttempts = retry
		}
	}
	
	// Set default timeout and retry attempts if not configured
	if cfg.Digiflazz.Timeout == 0 {
		cfg.Digiflazz.Timeout = 30 * time.Second
	}
	if cfg.Digiflazz.RetryAttempts == 0 {
		cfg.Digiflazz.RetryAttempts = 3
	}

	// Database configuration
	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Database.Port = p
		}
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		cfg.Database.Name = name
	}
	if user := os.Getenv("DB_USER"); user != "" {
		cfg.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		cfg.Database.Password = password
	}

	// Redis configuration
	if host := os.Getenv("REDIS_HOST"); host != "" {
		cfg.Redis.Host = host
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Redis.Port = p
		}
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		cfg.Redis.Password = password
	}

	// Logging configuration
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		cfg.Logging.Level = level
	}

	// Security configuration
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		cfg.Security.JWTSecret = jwtSecret
	}
	if rateLimit := os.Getenv("API_RATE_LIMIT"); rateLimit != "" {
		if rl, err := strconv.Atoi(rateLimit); err == nil {
			cfg.Security.APIRateLimit = rl
		}
	}

	// Monitoring configuration
	if enableMetrics := os.Getenv("ENABLE_METRICS"); enableMetrics != "" {
		cfg.Monitoring.EnableMetrics = enableMetrics == "true"
	}
	if metricsPort := os.Getenv("METRICS_PORT"); metricsPort != "" {
		if mp, err := strconv.Atoi(metricsPort); err == nil {
			cfg.Monitoring.MetricsPort = mp
		}
	}
}
