package bookmark

import (
	"context"

	"github.com/ThienKim52/golang-dev/internal/app/model"
	"github.com/ThienKim52/golang-dev/pkg/dbutils"
)

func (r *bookmarkRepo) CreateBookmark(ctx context.Context, bookmark *model.Bookmark) (*model.Bookmark, error) {
	err := r.db.WithContext(ctx).Create(bookmark).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}
	return bookmark, nil
}

func (r *bookmarkRepo) GetBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.Bookmark, int64, error) {
	bookmarks, err := r.getBookmarks(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	count, err := r.countBookmarks(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return bookmarks, count, nil
}

func (r *bookmarkRepo) getBookmarks(ctx context.Context, userId string, limit, offset int) ([]*model.Bookmark, error) {
	var bookmarks []*model.Bookmark
	err := r.db.WithContext(ctx).Where("user_id = ?", userId).Order("created_at ASC").Offset(offset).Limit(limit).Find(&bookmarks).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}
	return bookmarks, nil
}

func (r *bookmarkRepo) countBookmarks(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Bookmark{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, dbutils.CatchDBError(err)
	}
	return count, nil
}
