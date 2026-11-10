package user

import (
	"context"
	"github.com/ThienKim52/golang-dev/internal/app/model"
	"github.com/ThienKim52/golang-dev/pkg/dbutils"
)

func (r *sqlRepository) UpdateUser(ctx context.Context, userID string, updates map[string]interface{}) error  {
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		return dbutils.CatchDBError(err)
	}
	return nil
}
