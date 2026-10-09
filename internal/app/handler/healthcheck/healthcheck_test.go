package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ThienKim52/golang-dev/internal/app/service/healthcheck"
	mocks_healthcheck "github.com/ThienKim52/golang-dev/internal/app/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		serviceName  string
		instanceID   string
		setupMock    func(c *gin.Context) *mocks_healthcheck.HealthCheckService
		expectedCode int
	}{
		{
			name: "Successful health check",
			setupMock: func(c *gin.Context) *mocks_healthcheck.HealthCheckService {
				mockService := mocks_healthcheck.NewHealthCheckService(t)
				mockService.On("GetHealthCheck", c).Return(healthcheck.HealthCheckResult{
					Message:     "OK",
					ServiceName: "test-service",
					InstanceID:  "123e4567-e89b-12d3-a456-426614174000",
				}, nil)
				return mockService
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Health check with empty values",
			setupMock: func(c *gin.Context) *mocks_healthcheck.HealthCheckService {
				mockService := mocks_healthcheck.NewHealthCheckService(t)
				mockService.On("GetHealthCheck", c).Return(healthcheck.HealthCheckResult{
					Message:     "OK",
					ServiceName: "test-service",
					InstanceID:  "123e4567-e89b-12d3-a456-426614174000",
				}, nil)
				return mockService
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Redis connection error",
			setupMock: func(c *gin.Context) *mocks_healthcheck.HealthCheckService {
				mockService := mocks_healthcheck.NewHealthCheckService(t)
				mockService.On("GetHealthCheck", c).Return(healthcheck.HealthCheckResult{}, assert.AnError)
				return mockService
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create gin context
			w := httptest.NewRecorder()
            c, _ := gin.CreateTestContext(w)

			// init mock service for each case
			mockService := tt.setupMock(c)
			handler := NewHealthCheck(mockService)
			c.Request = httptest.NewRequest("GET", "/health-check", nil)

			// Call handler
			handler.GetResponse(c)

			// Assertions
			assert.Equal(t, tt.expectedCode, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}
