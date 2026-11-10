package bookmark

import (
	"context"

	"github.com/ThienKim52/golang-dev/internal/app/model"
	"github.com/ThienKim52/golang-dev/internal/app/repository/bookmark"
	"github.com/ThienKim52/golang-dev/pkg/utils"
)

type Bookmark interface {
	AddBookmark(ctx context.Context, description, url, userID string) (*model.Bookmark, error)
	GetBookmarks(ctx context.Context, userID string, page, limit int) (*GetBookmarksResult, error)
	UpdateBookmark(ctx context.Context, userID, bookmarkID, description, url string) error
	DeleteBookmark(ctx context.Context, userID, bookmarkID string) error
}

type BookmarkService struct {
	repo bookmark.BookmarkRepository
	keyGen utils.KeyGenerator
}

func NewBookmarkService(repo bookmark.BookmarkRepository, keyGen utils.KeyGenerator) Bookmark {
	return &BookmarkService{repo: repo, keyGen: keyGen}
}