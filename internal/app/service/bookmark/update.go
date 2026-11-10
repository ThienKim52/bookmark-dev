package bookmark

import (
	"context"
	"errors"
	"strings"
)

func (b *BookmarkService) UpdateBookmark(ctx context.Context, userID, bookmarkID, description, url string) error {
	if strings.TrimSpace(description) == "" ||
        strings.TrimSpace(url) == "" {
        return errors.New("description and url are required")
    }
	updates := map[string]interface{}{
        "description": description,
        "url":         url,
    }
	return b.repo.UpdateBookmark(
        ctx,
        userID,
        bookmarkID,
        updates,
    )
}
