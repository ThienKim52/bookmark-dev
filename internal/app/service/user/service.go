//go:generate go run github.com/vektra/mockery/v2@latest --name=Service --output=./mocks --outpkg=mocks_userservice
package user

import (
	"context"

	"github.com/ThienKim52/golang-dev/internal/app/model"
	"github.com/ThienKim52/golang-dev/internal/app/repository/user"
	"github.com/ThienKim52/golang-dev/pkg/jwtutils"
	"github.com/ThienKim52/golang-dev/pkg/utils"
)

type Service interface {
	CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error)
	Login(ctx context.Context, username, password string) (string, error)
	GetUserInfo(ctx context.Context, uid string) (*model.User, error)
	EditCurrentUser(ctx context.Context, userID, displayName, email string) error
}

type service struct {
	repo user.Repository
	hasher utils.Hasher
	jwtGenerator jwtutils.JWTGenerator
}

func NewService(repo user.Repository, hasher utils.Hasher, jwtGen jwtutils.JWTGenerator) Service{
	return &service{
		repo: repo,
		hasher: hasher,
		jwtGenerator: jwtGen,
	}
}