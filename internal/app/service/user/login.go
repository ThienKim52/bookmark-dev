package user

import (
	"context"
	"errors"
	"time"

	"github.com/ThienKim52/golang-dev/pkg/dbutils"
	"github.com/golang-jwt/jwt/v5"
)

const tokenDuration = 24 * time.Hour

var ErrInvalidCredentials = errors.New("invalid credentials")

func (s *service) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	switch {
	case errors.Is(err, dbutils.ErrRecordNotFound):
		return "", ErrInvalidCredentials

	case err == nil:
	default:
		return "", err
	}
	if !s.hasher.Compare(user.Password, password) {
		return "", ErrInvalidCredentials
	}
	tokenContent := jwt.MapClaims{
		"sub": user.ID,
		"email": user.Email,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(tokenDuration).Unix(),
	}
	tokenString, err := s.jwtGenerator.GenerateJWT(tokenContent)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}