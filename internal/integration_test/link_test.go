package integration_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"

	"testing"
	"github.com/ThienKim52/golang-dev/internal/api"
	redisPkg "github.com/ThienKim52/golang-dev/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestShortenURL(t *testing.T) {
	testCases := []struct {
		name                 string
		expectedResponseBody string
		reqBody              string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder
		setupRedis    func(ctx context.Context) *redis.Client

		expectedStatusCode int
	}{
		{
			name: "normal case",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest("POST", "/v1/links/shorten", bytes.NewBuffer([]byte(`{"url": "https://example.com", "exp": 604800}`)))
				respRecorder := httptest.NewRecorder()
				api.ServeHTTP(respRecorder, req)
				return respRecorder
			},
			setupRedis: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "Shorten URL generated successfully",
		},
		{
			name: "failed case",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest("POST", "/v1/links/shorten", bytes.NewBuffer([]byte(`{"url": "https://example.com", "exp": 604800}`)))
				respRecorder := httptest.NewRecorder()
				api.ServeHTTP(respRecorder, req)
				return respRecorder
			},
			setupRedis: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				_ = mock.Close()
				return mock
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"Processing error"}`,
		},
		{
			name: "normal case redirect",
			//setup request
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect/123456", nil)
				respRecorder := httptest.NewRecorder()
				api.ServeHTTP(respRecorder, req)
				return respRecorder
			},
			//setup redis
			setupRedis: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				err := mock.Set(ctx, "123456", "https://example.com", 0).Err()
				if err != nil {
					t.Fatalf("failed to set mock redis key: %v", err)
				}
				return mock
			},
			expectedStatusCode:   http.StatusFound,
			expectedResponseBody: "https://example.com",
		},
		{
			name: "failed case redirect",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect/123456", nil)
				respRecorder := httptest.NewRecorder()
				api.ServeHTTP(respRecorder, req)
				return respRecorder
			},
			setupRedis: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				_ = mock.Close()
				return mock
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"Processing error"}`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			// generate test redis
			mockRedis := tc.setupRedis(ctx)

			// Initialize the Gin router
			testRouter := api.New(&api.EngineOpts{
				Config: &api.Config{},
				RedisClient: mockRedis,
			})
			rec := tc.setupTestHTTP(testRouter)
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.expectedResponseBody)
		})
	}
}
