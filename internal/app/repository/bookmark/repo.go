package bookmark

import (
	"context"

	"github.com/ThienKim52/golang-dev/internal/app/model"
	"gorm.io/gorm"
)

type BookmarkRepository interface {
	CreateBookmark(ctx context.Context, bookmark *model.Bookmark) (*model.Bookmark, error) 
	UpdateBookmark(ctx context.Context, userID, bookmarkID string, updates map[string]interface{}) error
	GetBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.Bookmark, int64, error)
	DeleteBookmark(ctx context.Context, userID, bookmarkID string) error
	CodeExists(ctx context.Context, code string) (bool, error)
}

type bookmarkRepo struct {
	db *gorm.DB
}

func NewBookmarkRepository(db *gorm.DB) BookmarkRepository {
	return &bookmarkRepo{db: db}
}

func (r *bookmarkRepo) CodeExists(ctx context.Context, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Bookmark{}).Where("code = ?", code).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}	
