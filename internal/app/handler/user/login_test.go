package user

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ThienKim52/golang-dev/internal/app/service/user"
	mocks_userservice "github.com/ThienKim52/golang-dev/internal/app/service/user/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testCases := []struct {
		name               string
		reqBody            string
		setupMock          func(c *gin.Context) *mocks_userservice.Service
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:    "Login successfully",
			reqBody: `{"username":"test123","password":"12223344"}`,
			setupMock: func(c *gin.Context) *mocks_userservice.Service {
				mock := mocks_userservice.NewService(t)
				mock.On("Login", c, "test123", "12223344").Return("fake-jwt-token", nil)
				return mock
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `"Logged in successfully!"`,
		},
		{
			name:     "Invalid credentials",
			reqBody: `{"username":"test123","password":"12223344"}`,
			setupMock: func(c *gin.Context) *mocks_userservice.Service {
				mock := mocks_userservice.NewService(t)
				mock.On("Login", c, "test123", "wrongpass").Return("", user.ErrInvalidCredentials)
				return mock
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
            name:    "Password too short (< 8)",
            reqBody: `{"username":"test","password":"short"}`,
            setupMock: func(c *gin.Context) *mocks_userservice.Service {
                return mocks_userservice.NewService(t) // won't be called
            },
            expectedStatusCode: http.StatusBadRequest,
        },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			mockService := tc.setupMock(c)
			hdl := NewHandler(mockService)
			c.Request = httptest.NewRequest("POST", "/v1/users/login", bytes.NewBufferString(tc.reqBody))
			c.Request.Header.Set("Content-Type", "application/json")
			hdl.Login(c)
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			if tc.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tc.expectedBody)
			}
		})
	}
}
