package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ThienKim52/golang-dev/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ThienKim52/golang-dev/pkg/sqldb"
)

func TestHealthCheck(t *testing.T) {
	// not t.Parallel() because of os.Setenv
	t.Parallel()
	testCases := []struct {
		name               string
		cfg                api.Config
		expectedStatusCode int
		expectedMessage    string
		verifyResponse     func(t *testing.T, body string) // verify response
	}{
		{
			name: "With Set Env",
			cfg: api.Config{
				ServiceName: "test-service",
				InstanceID:  "test-instance-123",
				AppPort:     "8081",
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "OK",
			verifyResponse: func(t *testing.T, body string) {
				var resp map[string]string
				err := json.Unmarshal([]byte(body), &resp)
				require.NoError(t, err)
				assert.Equal(t, "test-service", resp["service_name"])
				assert.Equal(t, "test-instance-123", resp["instance_id"])
			},
		},
		{
			name: "With Generated UUID",
			cfg: api.Config{
				ServiceName: "test-service-uuid",
				InstanceID:  "test-instance-456",
				AppPort:     "8082",
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "OK",
			verifyResponse: func(t *testing.T, body string) {
				var resp map[string]string
				err := json.Unmarshal([]byte(body), &resp)
				require.NoError(t, err)
				assert.Equal(t, "test-service-uuid", resp["service_name"])
				assert.Equal(t, "test-instance-456", resp["instance_id"])
			},
		},
		{
			name: "With default env",
			cfg: api.Config{
				ServiceName: "health-check-service",
				InstanceID:  "test-instance-123",
				AppPort:     "8082",
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "OK",
			verifyResponse: func(t *testing.T, body string) {
				var resp map[string]string
				err := json.Unmarshal([]byte(body), &resp)
				require.NoError(t, err)
				assert.Equal(t, "health-check-service", resp["service_name"])
				assert.Equal(t, "test-instance-123", resp["instance_id"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// create virtual database miniredis
			mr, err := miniredis.Run()
			require.NoError(t, err)
			defer mr.Close()

			// create client connected to miniredis
			redisClient := redis.NewClient(&redis.Options{
				Addr: mr.Addr(),
			})
			defer redisClient.Close()
			db := sqldb.InitMockDB(t)

			// create engine application
			apiEngine := api.New(&api.EngineOpts{
				App:         gin.New(),
				Config:      &tc.cfg,
				RedisClient: redisClient,
				DB:          db,
				
			})

			// simulate HTTP request to endpoint /health-check
			req := httptest.NewRequest("GET", "/health-check", nil)
			respRec := httptest.NewRecorder()
			apiEngine.ServeHTTP(respRec, req)

			// assert response
			assert.Equal(t, tc.expectedStatusCode, respRec.Code)
			assert.Contains(t, respRec.Body.String(), tc.expectedMessage)

			// verify response JSON
			if tc.verifyResponse != nil {
				tc.verifyResponse(t, respRec.Body.String())
			}
		})
	}
}
