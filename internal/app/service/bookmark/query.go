package bookmark
import
(
	"context"
	"github.com/ThienKim52/golang-dev/internal/app/model"
)

type GetBookmarksResult struct {
	Bookmarks []*model.Bookmark `json:"bookmarks"`
	Total int64 `json:"total"`
}

func (s *BookmarkService) GetBookmarks(ctx context.Context, userID string, page, limit int) (*GetBookmarksResult, error) {
	offset := (page - 1) * limit
	bookmarks, count, err := s.repo.GetBookmarks(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	return &GetBookmarksResult{
		Bookmarks: bookmarks,
		Total: count,
	}, nil
}