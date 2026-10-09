package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	mocks_genpass "github.com/ThienKim52/golang-dev/internal/app/service/mocks"
)

var testErr = errors.New("test error")

func TestGenPass_GeneratePassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockSvc     func() *mocks_genpass.GenPass
		setupTestRequest func(c *gin.Context)

		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "normal case",

			setupMockSvc: func() *mocks_genpass.GenPass {
				mockSvc := mocks_genpass.NewGenPass(t)
				mockSvc.On("GeneratePassword", passwordLength).Return("123456789012", nil)
				return mockSvc
			},
			setupTestRequest: func(c *gin.Context) {
				c.Request = httptest.NewRequest(http.MethodGet, "/genpass", nil)
			},

			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"password":"123456789012"}`,
		},
		{
			name: "err case",

			setupMockSvc: func() *mocks_genpass.GenPass {
				mockSvc := mocks_genpass.NewGenPass(t)
				mockSvc.On("GeneratePassword", passwordLength).Return("", testErr)
				return mockSvc
			},
			setupTestRequest: func(c *gin.Context) {
				c.Request = httptest.NewRequest(http.MethodGet, "/genpass", nil)
			},

			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"message":"Processing error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)

			tc.setupTestRequest(c)

			mockSvc := tc.setupMockSvc()
			genPassHandler := NewGenPass(mockSvc)
			genPassHandler.GeneratePassword(c)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Body.String())

		})
	}
}
