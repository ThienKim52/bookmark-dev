package bookmark

import (
	"github.com/ThienKim52/golang-dev/internal/app/service/bookmark"
	"github.com/gin-gonic/gin"
)

type BookmarkHandler interface {
	AddBookmark(c *gin.Context)
	GetBookmarks(c *gin.Context)
	UpdateBookmark(c *gin.Context)
	DeleteBookmark(c *gin.Context)
}

type bookmarkHandler struct {
	bookmarkSvc bookmark.Bookmark
}

func NewBookmarkHandler(bookmarkSvc bookmark.Bookmark) BookmarkHandler {
	return &bookmarkHandler{bookmarkSvc: bookmarkSvc}
}
