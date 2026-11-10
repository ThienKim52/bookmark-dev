package bookmark

import (
	"context"

	"github.com/ThienKim52/golang-dev/internal/app/model"
)

const codeLength = 8

func (b *BookmarkService) AddBookmark(ctx context.Context, description, url, userID string) (*model.Bookmark, error) {
	// Generate unique code - check both Redis and DB
	var key string
	var err error
	
	// Try to generate a unique code
	for {
		key, err = b.keyGen.KeyGenerator(codeLength)
		if err != nil {
			return nil, err
		}
		
		// Check if code exists in DB
		exists, err := b.repo.CodeExists(ctx, key)
		if err != nil {
			return nil, err
		}
		if exists {
			continue // Code exists, generate a new one
		}
		
		// Code is unique, break the loop
		break
	}
	
	bookmark := &model.Bookmark{
		Description: description,
		URL:         url,
		Code:        key,
		UserID:      userID,
	}
	res, err := b.repo.CreateBookmark(ctx, bookmark)
	return res, err
}