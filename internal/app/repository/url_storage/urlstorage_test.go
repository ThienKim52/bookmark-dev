package urlstorage

import (
	"context"
	"testing"

	redis2 "github.com/ThienKim52/golang-dev/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestUrlStorage_StoreURL(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name        string
		setupMock   func(ctx context.Context, t *testing.T) *redis.Client
		expectedErr error
		verifyFunc  func(ctx context.Context, r *redis.Client)
	}{
		{
			name: "normal",
			setupMock: func(ctx context.Context, t *testing.T) *redis.Client {
				mock := redis2.InitMockRedis(t)
				return mock
			},
			expectedErr: nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				res, err := r.Get(ctx, "test").Result()
				assert.NoError(t, err)
				assert.Equal(t, "http://google.com", res)
			},
		},
		{
			name: "connection error",
			setupMock: func(ctx context.Context, t *testing.T) *redis.Client {
				mock := redis2.InitMockRedis(t)
				_ = mock.Close()
				return mock
			},
			expectedErr: redis.ErrClosed,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mock := tc.setupMock(ctx, t)
			storage := NewURLStorage(mock)
			err := storage.StoreURL(ctx, "test", "http://google.com", 0)
			assert.Equal(t, tc.expectedErr, err)

			if tc.verifyFunc != nil {
				tc.verifyFunc(ctx, mock)
			}
		})
	}
}

func TestUrlStorage_GetURL(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name        string
		setupMock   func(ctx context.Context, t *testing.T) *redis.Client
		expectedErr error
		expectedURL string
		verifyFunc  func(ctx context.Context, r *redis.Client)
	}{
		{
			name: "normal",
			setupMock: func(ctx context.Context, t *testing.T) *redis.Client {
				mock := redis2.InitMockRedis(t)
				return mock
			},
			expectedErr: nil,
			expectedURL: "http://google.com",
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				_, err := r.Set(ctx, "test", "http://google.com", 0).Result()
				assert.NoError(t, err)
			},
		},
		{
			name: "connection error",
			setupMock: func(ctx context.Context, t *testing.T) *redis.Client {
				mock := redis2.InitMockRedis(t)
				_ = mock.Close()
				return mock
			},
			expectedErr: redis.ErrClosed,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mock := tc.setupMock(ctx, t)
			storage := NewURLStorage(mock)
			if tc.verifyFunc != nil {
				tc.verifyFunc(ctx, mock)
			}
			url, err := storage.GetURL(ctx, "test")
			assert.Equal(t, tc.expectedErr, err)
			if url != "" {
				assert.Equal(t, tc.expectedURL, url)
			}
		})
	}
}
