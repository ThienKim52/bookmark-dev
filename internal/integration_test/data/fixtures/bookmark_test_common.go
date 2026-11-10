package fixtures

import (
	"github.com/ThienKim52/golang-dev/internal/app/model"
	"gorm.io/gorm"
)

type BookmarkCommonTestDB struct {
	UserCommonTestDB
}

func (b *BookmarkCommonTestDB) Migrate() error {
	return b.db.AutoMigrate(&model.User{}, &model.Bookmark{})
}

func (b *BookmarkCommonTestDB) GenerateData() error {
	err := b.UserCommonTestDB.GenerateData()
	if err != nil {
		return err
	}
	bookmarks := []*model.Bookmark{
		{
			Base: GetTestBase("deb745af-1a62-4efa-99a0-f06b274bd990"),
			Description: "Bookmark 1",
			URL: "https://example.com/1",
			Code: "bookmark-1",
			UserID: "user-1",
		},
		{
			Base: GetTestBase("deb745af-1a62-4efa-99a0-f06b274bd991"),
			Description: "Bookmark 2",
			URL: "https://example.com/2",
			Code: "bookmark-2",
			UserID: "user-2",
		},
	}
	return b.db.Session(&gorm.Session{SkipHooks: true}).CreateInBatches(bookmarks, 10).Error
}