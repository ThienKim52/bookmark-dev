package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ThienKim52/golang-dev/internal/app/model"
	mocks_userservice "github.com/ThienKim52/golang-dev/internal/app/service/user/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestSelfInfo(t *testing.T) {
	testCases := []struct {
		name      string
		setupMock func(c *gin.Context) *mocks_userservice.Service
		setupTestRequest func(c *gin.Context)
		expectedStatusCode int
	}{
		{
			name: "Get info successfully",
			setupMock: func(c *gin.Context) *mocks_userservice.Service {
				mock := mocks_userservice.NewService(t)
				mock.On("GetUserInfo", c, "test-user-id").Return(&model.User{
					ID: "test-user-id",
					Username: "test",
					Email: "test@gmail.com",
					DisplayName: "test123",
				}, nil)
				return mock
			},
			setupTestRequest: func(c *gin.Context) {
				c.Request = httptest.NewRequest("GET", "/v1/self/info", nil)
				c.Set("claims", jwt.MapClaims{
					"sub": "test-user-id",
					"email": "test@gmail.com",
				})
			},
			expectedStatusCode: http.StatusOK,
		},
		{
            name: "No claims in context → 401",
            setupMock: func(c *gin.Context) *mocks_userservice.Service {
                return mocks_userservice.NewService(t)
            },
            setupTestRequest: func(c *gin.Context) {
                c.Request = httptest.NewRequest("GET", "/v1/self/info", nil)
                // Không set claims → mong đợi 401
            },
            expectedStatusCode: http.StatusUnauthorized,
        },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			mockService := tc.setupMock(c)
			tc.setupTestRequest(c)
			hdl := NewHandler(mockService)
			hdl.GetUserInfo(c)
			assert.Equal(t, tc.expectedStatusCode, w.Code)
		})
	}
		
}