package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// LinkRepository defines the interface for link storage operations
type LinkRepository interface {
	StoreURL(ctx context.Context, code, url string, exp time.Duration) error
	GetURL(ctx context.Context, code string) (string, error)
	Exists(ctx context.Context, code string) (bool, error)
	Ping(ctx context.Context) error
}

// RedisLinkRepository implements LinkRepository using Redis
type RedisLinkRepository struct {
	c *redis.Client
}

// NewRedisLinkRepository creates a new Redis link repository
func NewRedisLinkRepository(c *redis.Client) LinkRepository {
	return &RedisLinkRepository{c: c}
}

// Save saves a link to Redis with expiration
func (r *RedisLinkRepository) StoreURL(ctx context.Context, code, url string, exp time.Duration) error {
	return r.c.Set(ctx, code, url, exp).Err()
}

// Get retrieves a link by its code
func (r *RedisLinkRepository) GetURL(ctx context.Context, code string) (string, error) {
	return r.c.Get(ctx, code).Result()
}

// Exists checks if a code already exists in Redis
func (r *RedisLinkRepository) Exists(ctx context.Context, code string) (bool, error) {
	exists, err := r.c.Exists(ctx, code).Result()
	return exists > 0, err
}

// Ping checks Redis connection
func (r *RedisLinkRepository) Ping(ctx context.Context) error {
	return r.c.Ping(ctx).Err()
}
