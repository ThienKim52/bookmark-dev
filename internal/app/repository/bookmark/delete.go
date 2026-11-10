package bookmark

import (
	"context"

	"github.com/ThienKim52/golang-dev/internal/app/model"
	"github.com/ThienKim52/golang-dev/pkg/dbutils"
)

func (r *bookmarkRepo) DeleteBookmark(ctx context.Context, userID, bookmarkID string) error {
	result  := r.db.WithContext(ctx).Model(&model.Bookmark{}).Where("id = ? AND user_id = ?", bookmarkID, userID).Delete(&model.Bookmark{})
	if result.Error != nil {
        return dbutils.CatchDBError(result.Error)
    }
	if result.RowsAffected == 0 {
		return dbutils.ErrBookmarkNotFound
	}	
    return nil
}
