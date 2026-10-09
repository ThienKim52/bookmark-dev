package user

import (
	"context"

	"github.com/ThienKim52/golang-dev/internal/app/model"
)

func (s *service) GetUserInfo(ctx context.Context, uid string) (*model.User, error) {
	return s.repo.GetUserByID(ctx, uid)
}