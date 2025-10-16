package tests

import (
	"testing"

	"gateway-digiflazz/internal/config"
	"gateway-digiflazz/pkg/digiflazz"

	"github.com/stretchr/testify/assert"
)

func TestDigiflazzClient(t *testing.T) {
	// Load test configuration
	cfg, err := config.Load()
	assert.NoError(t, err)

	// Create client
	client := digiflazz.NewClient(cfg.Digiflazz)

	// Test balance check (this will fail without valid credentials)
	t.Run("CheckBalance", func(t *testing.T) {
		// Skip if no credentials provided
		if cfg.Digiflazz.Username == "" || cfg.Digiflazz.APIKey == "" {
			t.Skip("Skipping test: No Digiflazz credentials provided")
		}

		balance, err := client.CheckBalance()
		assert.NoError(t, err)
		assert.NotNil(t, balance)
	})

	// Test price list (this will fail without valid credentials)
	t.Run("GetPrices", func(t *testing.T) {
		// Skip if no credentials provided
		if cfg.Digiflazz.Username == "" || cfg.Digiflazz.APIKey == "" {
			t.Skip("Skipping test: No Digiflazz credentials provided")
		}

		prices, err := client.GetPrices("prabayar")
		assert.NoError(t, err)
		assert.NotNil(t, prices)
	})
}

func TestConfig(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.Server.Port)
	assert.NotEmpty(t, cfg.Digiflazz.BaseURL)
}
