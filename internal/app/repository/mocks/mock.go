package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockLinkRepository is a mock implementation of LinkRepository
type MockLinkRepository struct {
	mock.Mock
}

func (m *MockLinkRepository) StoreURL(ctx context.Context, code, url string, exp time.Duration) error {
	args := m.Called(ctx, code, url, exp)
	return args.Error(0)
}

func (m *MockLinkRepository) GetURL(ctx context.Context, code string) (string, error) {
	args := m.Called(ctx, code)
	return args.String(0), args.Error(1)
}

func (m *MockLinkRepository) Exists(ctx context.Context, code string) (bool, error) {
	args := m.Called(ctx, code)
	return args.Bool(0), args.Error(1)
}

func (m *MockLinkRepository) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
