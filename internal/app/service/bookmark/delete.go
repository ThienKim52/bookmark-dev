package bookmark

import (
	"context"
)

func (b *BookmarkService) DeleteBookmark(ctx context.Context, userID, bookmarkID string) error {
	return b.repo.DeleteBookmark(ctx, userID, bookmarkID)
}
