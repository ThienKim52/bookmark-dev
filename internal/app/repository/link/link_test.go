package repository

import (
	"context"
	"testing"
	"time"

	redis2 "github.com/ThienKim52/golang-dev/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupMockRedis(t *testing.T) *redis.Client {
			mock := redis2.InitMockRedis(t)
			return mock
		}

func TestRedisLinkRepository(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name        string
		expectedErr error
		verifyFunc  func(ctx context.Context, repo *LinkRepository, r *redis.Client)
	}{
		{
			name: "new Redis LinkRepository",
			expectedErr: nil,
			verifyFunc: func(ctx context.Context, repo *LinkRepository, r *redis.Client) {
				assert.NotNil(t, repo)
				_, ok := (*repo).(*RedisLinkRepository)
				assert.True(t, ok)
			},
		},
		{
			name: "func StoreURL",
			verifyFunc: func(ctx context.Context, repo *LinkRepository, r *redis.Client) {
				err := (*repo).StoreURL(ctx, "test", "http://google.com", 0)
				assert.NoError(t, err)
				val := r.Get(ctx, "test")
				assert.NotNil(t, val)
				assert.Equal(t, "http://google.com", val.Val())
			},
			expectedErr: redis.ErrClosed,
		},
		{
			name: "func GetURL",
			verifyFunc: func(ctx context.Context, repo *LinkRepository, r *redis.Client) {
				r.Set(ctx, "test", "http://google.com", time.Hour)
				val, err := (*repo).GetURL(ctx, "test")
				assert.NoError(t, err)
				assert.Equal(t, "http://google.com", val)

				val, err = (*repo).GetURL(ctx, "non-existent")
				assert.ErrorIs(t, err, redis.Nil)
				assert.Empty(t, val)
			},
			expectedErr: redis.ErrClosed,
		},
		{
			name: "func Exists",
			verifyFunc: func(ctx context.Context, repo *LinkRepository, r *redis.Client) {
				exists, err := (*repo).Exists(ctx, "test")
				assert.NoError(t, err)
				assert.False(t, exists)

				r.Set(ctx, "test", "http://google.com", time.Hour)

				exists, err = (*repo).Exists(ctx, "test")
				assert.NoError(t, err)
				assert.True(t, exists)
			},
			expectedErr: redis.ErrClosed,
		},
		{
			name: "func Ping",
			verifyFunc: func(ctx context.Context, repo *LinkRepository, r *redis.Client){
				err := (*repo).Ping(ctx)
				assert.NoError(t, err)
			},
			expectedErr: redis.ErrClosed,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mock := setupMockRedis(t)
			repo := NewRedisLinkRepository(mock)
			if tc.verifyFunc != nil {
				tc.verifyFunc(ctx, &repo, mock)
			}
		})
	}
}
