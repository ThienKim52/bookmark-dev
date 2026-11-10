package link

import (
	"context"
	"errors"
	"time"

	"github.com/ThienKim52/golang-dev/internal/app/repository/link"
	"github.com/ThienKim52/golang-dev/internal/app/service/genpass"
)

const (
	codeLength = 7
	maxRetries = 10
)

var ErrMaxRetriesExceeded = errors.New("max retries exceeded while generating unique code")

//go:generate mockery --name=LinkService --filename=linkservice.go --outpkg=mocks_genpass
type LinkService interface {
	ShortenURL(ctx context.Context, url string, exp time.Duration) (code string, err error)
	GetLinkFromCode(ctx context.Context, code string) (string, error)
}

// linkService implements LinkService
type linkService struct {
	storage repository.LinkRepository
	genPass genpass.GenPass
}

// NewLinkService creates a new link service
func NewLinkService(storage repository.LinkRepository, genPass genpass.GenPass) LinkService {
	return &linkService{
		storage: storage,
		genPass: genPass,
	}
}

// ShortenURL creates a short code for a URL and saves it to the repository
func (s *linkService) ShortenURL(ctx context.Context, url string, exp time.Duration) (string, error) {
	var code string
	var err error

	// Generate a unique code with retry mechanism
	for i := 0; i < maxRetries; i++ {
		code, err = s.genPass.GeneratePassword(codeLength)
		if err != nil {
			return "", err
		}

		// Check if code already exists
		exists, err := s.storage.Exists(ctx, code)
		if err != nil {
			return "", err
		}

		if !exists {
			// Code is unique, save the link
			if err := s.storage.StoreURL(ctx, code, url, exp); err != nil {
				return "", err
			}

			return code, nil
		}
		// Code exists, try again
	}

	return "", ErrMaxRetriesExceeded
}

func (s *linkService) GetLinkFromCode(ctx context.Context, code string) (string, error) {
	return s.storage.GetURL(ctx, code)
}
