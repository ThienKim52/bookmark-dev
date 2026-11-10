package user

import (
	"context"

	"github.com/ThienKim52/golang-dev/internal/app/model"
	"github.com/ThienKim52/golang-dev/pkg/dbutils"
)

func (r *sqlRepository) CreateUser(ctx context.Context, newUser *model.User) (*model.User, error) {
	db := r.db.WithContext(ctx)

	result := db.Create(newUser)

	err := result.Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}
	return newUser, nil
}
