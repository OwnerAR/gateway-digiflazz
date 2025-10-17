package cache

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteCache implements CacheInterface using SQLite
type SQLiteCache struct {
	db *sql.DB
}

// NewSQLiteCache creates a new SQLite cache instance
func NewSQLiteCache(dbPath string) (*SQLiteCache, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	cache := &SQLiteCache{db: db}
	
	// Create cache table if not exists
	if err := cache.createTable(); err != nil {
		return nil, err
	}

	return cache, nil
}

// createTable creates the cache table
func (s *SQLiteCache) createTable() error {
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
func (s *SQLiteCache) Get(ctx context.Context, key string) (string, error) {
	// Handle permanent cache (TTL = 0) by checking expires_at with timezone support
	query := `SELECT data FROM pln_inquiry_cache WHERE customer_no = ? AND (expires_at > ? OR expires_at LIKE '0001-01-01%')`
	
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
func (s *SQLiteCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	now := time.Now()
	var expiresAt time.Time
	
	if ttl == 0 {
		// Permanent cache - use zero time
		expiresAt = time.Time{}
	} else {
		expiresAt = now.Add(ttl)
	}
	
	query := `
	INSERT OR REPLACE INTO pln_inquiry_cache (customer_no, data, created_at, expires_at)
	VALUES (?, ?, ?, ?)
	`
	
	_, err := s.db.ExecContext(ctx, query, key, value, now, expiresAt)
	return err
}

// Delete removes a value from cache
func (s *SQLiteCache) Delete(ctx context.Context, key string) error {
	query := `DELETE FROM pln_inquiry_cache WHERE customer_no = ?`
	_, err := s.db.ExecContext(ctx, query, key)
	return err
}

// DeleteExpired removes expired entries from cache
func (s *SQLiteCache) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM pln_inquiry_cache WHERE expires_at <= ?`
	_, err := s.db.ExecContext(ctx, query, time.Now())
	return err
}

// ClearAll removes all entries from cache
func (s *SQLiteCache) ClearAll(ctx context.Context) error {
	query := `DELETE FROM pln_inquiry_cache`
	_, err := s.db.ExecContext(ctx, query)
	return err
}

// GetStats returns cache statistics
func (s *SQLiteCache) GetStats(ctx context.Context) (map[string]interface{}, error) {
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
func (s *SQLiteCache) Close() error {
	return s.db.Close()
}

// Ping tests the connection to SQLite
func (s *SQLiteCache) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}
