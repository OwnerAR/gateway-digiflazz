package cache

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// Pure Go SQLite driver - no CGO required
	_ "github.com/crawshaw/sqlite/sqlite3"
)

// SQLitePureCache implements CacheInterface using pure Go SQLite
// This version doesn't require CGO, making cross-compilation easier
type SQLitePureCache struct {
	db *sql.DB
}

// NewSQLitePureCache creates a new pure Go SQLite cache instance
func NewSQLitePureCache(dbPath string) (*SQLitePureCache, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	cache := &SQLitePureCache{db: db}
	
	// Create cache table if not exists
	if err := cache.createTable(); err != nil {
		return nil, err
	}

	return cache, nil
}

// createTable creates the cache table
func (s *SQLitePureCache) createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS pln_inquiry_cache (
		customer_no TEXT PRIMARY KEY,
		data TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		expires_at DATETIME NOT NULL
	);
	
	CREATE INDEX IF NOT EXISTS idx_expires_at ON pln_inquiry_cache(expires_at);
	`
	
	_, err := s.db.Exec(query)
	return err
}

// Get retrieves a value from cache
func (s *SQLitePureCache) Get(ctx context.Context, key string) (string, error) {
	query := `SELECT data FROM pln_inquiry_cache WHERE customer_no = ? AND expires_at > ?`
	
	var data string
	err := s.db.QueryRowContext(ctx, query, key, time.Now()).Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("key not found")
		}
		return "", err
	}
	
	return data, nil
}

// Set stores a value in cache with TTL
func (s *SQLitePureCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	now := time.Now()
	expiresAt := now.Add(ttl)
	
	query := `
	INSERT OR REPLACE INTO pln_inquiry_cache (customer_no, data, created_at, expires_at)
	VALUES (?, ?, ?, ?)
	`
	
	_, err := s.db.ExecContext(ctx, query, key, value, now, expiresAt)
	return err
}

// Delete removes a value from cache
func (s *SQLitePureCache) Delete(ctx context.Context, key string) error {
	query := `DELETE FROM pln_inquiry_cache WHERE customer_no = ?`
	_, err := s.db.ExecContext(ctx, query, key)
	return err
}

// DeleteExpired removes expired entries from cache
func (s *SQLitePureCache) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM pln_inquiry_cache WHERE expires_at <= ?`
	_, err := s.db.ExecContext(ctx, query, time.Now())
	return err
}

// ClearAll removes all entries from cache
func (s *SQLitePureCache) ClearAll(ctx context.Context) error {
	query := `DELETE FROM pln_inquiry_cache`
	_, err := s.db.ExecContext(ctx, query)
	return err
}

// GetStats returns cache statistics
func (s *SQLitePureCache) GetStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Total entries
	var total int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM pln_inquiry_cache`).Scan(&total)
	if err != nil {
		return nil, err
	}
	stats["total_entries"] = total
	
	// Expired entries
	var expired int
	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM pln_inquiry_cache WHERE expires_at <= ?`, time.Now()).Scan(&expired)
	if err != nil {
		return nil, err
	}
	stats["expired_entries"] = expired
	
	// Active entries
	stats["active_entries"] = total - expired
	
	return stats, nil
}

// Close closes the database connection
func (s *SQLitePureCache) Close() error {
	return s.db.Close()
}

// Ping tests the connection to SQLite
func (s *SQLitePureCache) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}
