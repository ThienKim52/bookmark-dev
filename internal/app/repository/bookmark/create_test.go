package bookmark

import (
	"testing"

	"github.com/ThienKim52/golang-dev/internal/app/model"
	"github.com/ThienKim52/golang-dev/internal/integration_test/data/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestBookmark(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		setupDB       func(t *testing.T) *gorm.DB
		input         *model.Bookmark
		expectedError error
		verifyFunc    func(db *gorm.DB, expected *model.Bookmark)
	}{
		{
			name: "Success",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{}).DB()
			},
			input: &model.Bookmark{
				Base:        fixtures.GetTestBase("deb745af-1a62-4efa-99a0-f06b274bd992"),
				Description: "Bookmark 3",
				URL:         "https://example.com/3",
				Code:        "bookmark-3",
				UserID:      "user-3",
			},
			expectedError: nil,
			verifyFunc: func(db *gorm.DB, expected *model.Bookmark) {
				var bookmark model.Bookmark
				err := db.Where("code = ?", "bookmark-3").First(&bookmark).Error
				require.NoError(t, err)
				assert.Equal(t, expected.ID, bookmark.ID)
				assert.Equal(t, expected.Code, bookmark.Code)
				assert.Equal(t, expected.Description, bookmark.Description)
				assert.Equal(t, expected.URL, bookmark.URL)
				assert.Equal(t, expected.UserID, bookmark.UserID)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := tc.setupDB(t)
			ctx := t.Context()
			repo := NewBookmarkRepository(db)
			output, err := repo.CreateBookmark(ctx, tc.input)
			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, output)
			}
			if tc.verifyFunc != nil {
				tc.verifyFunc(db, output)
			}
		})
	}
}