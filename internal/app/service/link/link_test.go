package link

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ThienKim52/golang-dev/internal/app/repository/mocks"
	urlstorage "github.com/ThienKim52/golang-dev/internal/app/repository/url_storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGenPass is a mock implementation of GenPass
type MockGenPass struct {
	mock.Mock
}

func (m *MockGenPass) GeneratePassword(passwordLength int) (string, error) {
	args := m.Called(passwordLength)
	return args.String(0), args.Error(1)
}

func TestNewLinkService(t *testing.T) {
	mockRepo := new(mocks.MockLinkRepository)
	mockGenPass := new(MockGenPass)
	service := NewLinkService(mockRepo, mockGenPass)
	assert.NotNil(t, service)
	_, ok := service.(*linkService)
	assert.True(t, ok)
}

func TestLinkService_ShortenURL(t *testing.T) {
	repoErr := errors.New("repository error")

	tests := []struct {
		name         string
		url          string
		exp          time.Duration
		setupMocks   func(mockRepo *mocks.MockLinkRepository, mockGenPass *MockGenPass)
		expectedCode string
		expectedErr  error
	}{
		{
			name: "Success",
			url:  "https://example.com",
			exp:  604800 * time.Second,
			setupMocks: func(mockRepo *mocks.MockLinkRepository, mockGenPass *MockGenPass) {
				ctx := context.Background()
				mockGenPass.On("GeneratePassword", codeLength).Return("abc1234", nil)
				mockRepo.On("Exists", ctx, "abc1234").Return(false, nil)
				mockRepo.On("StoreURL", ctx, "abc1234", "https://example.com", 604800*time.Second).Return(nil)
			},
			expectedCode: "abc1234",
			expectedErr:  nil,
		},
		{
			name: "CodeConflict_RetrySuccess",
			url:  "https://example.com",
			exp:  604800 * time.Second,
			setupMocks: func(mockRepo *mocks.MockLinkRepository, mockGenPass *MockGenPass) {
				ctx := context.Background()
				mockGenPass.On("GeneratePassword", codeLength).Return("code1", nil).Once()
				mockRepo.On("Exists", ctx, "code1").Return(true, nil).Once()
				mockGenPass.On("GeneratePassword", codeLength).Return("code2", nil).Once()
				mockRepo.On("Exists", ctx, "code2").Return(false, nil).Once()
				mockRepo.On("StoreURL", ctx, "code2", "https://example.com", 604800*time.Second).Return(nil)
			},
			expectedCode: "code2",
			expectedErr:  nil,
		},
		{
			name: "MaxRetriesExceeded",
			url:  "https://example.com",
			exp:  604800 * time.Second,
			setupMocks: func(mockRepo *mocks.MockLinkRepository, mockGenPass *MockGenPass) {
				ctx := context.Background()
				mockGenPass.On("GeneratePassword", codeLength).Return("code", nil)
				mockRepo.On("Exists", ctx, mock.AnythingOfType("string")).Return(true, nil)
			},
			expectedCode: "",
			expectedErr:  ErrMaxRetriesExceeded,
		},
		{
			name: "RepositoryError",
			url:  "https://example.com",
			exp:  604800 * time.Second,
			setupMocks: func(mockRepo *mocks.MockLinkRepository, mockGenPass *MockGenPass) {
				ctx := context.Background()
				mockGenPass.On("GeneratePassword", codeLength).Return("abc1234", nil)
				mockRepo.On("Exists", ctx, "abc1234").Return(false, repoErr)
			},
			expectedCode: "",
			expectedErr:  repoErr,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.MockLinkRepository)
			mockGenPass := new(MockGenPass)
			service := NewLinkService(mockRepo, mockGenPass)

			tc.setupMocks(mockRepo, mockGenPass)

			ctx := context.Background()
			code, err := service.ShortenURL(ctx, tc.url, tc.exp)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Empty(t, code)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCode, code)
			}

			mockRepo.AssertExpectations(t)
			mockGenPass.AssertExpectations(t)
		})
	}
}

func TestLinkService_GetLinkFromCode(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		setupMocks   func(mockRepo *mocks.MockLinkRepository)
		expectedLink string
		expectedErr  error
	}{
		{
			name: "Success",
			code: "abc1234",
			setupMocks: func(mockRepo *mocks.MockLinkRepository) {
				ctx := context.Background()
				mockRepo.On("GetURL", ctx, "abc1234").Return("https://example.com", nil)
			},
			expectedLink: "https://example.com",
			expectedErr:  nil,
		},
		{
			name: "CodeNotFound",
			code: "abc1234",
			setupMocks: func(mockRepo *mocks.MockLinkRepository) {
				ctx := context.Background()
				mockRepo.On("GetURL", ctx, "abc1234").Return("", urlstorage.ErrCodeNotFound)
			},
			expectedLink: "",
			expectedErr:  urlstorage.ErrCodeNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.MockLinkRepository)
			mockGenPass := new(MockGenPass)
			service := NewLinkService(mockRepo, mockGenPass)

			tc.setupMocks(mockRepo)

			ctx := context.Background()
			link, err := service.GetLinkFromCode(ctx, tc.code)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Empty(t, link)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedLink, link)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}