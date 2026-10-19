package user_test

import (
	"context"
	"testing"

	"github.com/ThienKim52/golang-dev/internal/app/model"
	"github.com/ThienKim52/golang-dev/internal/app/repository/user"
	"github.com/ThienKim52/golang-dev/pkg/sqldb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestGetUserInfo(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		inputUsername string 
		setupDatabase func(t *testing.T, db *gorm.DB) // setup dữ liệu
		expectedUser  *model.User
		expectedError bool
	}{
		{
			name:    "User found",
			inputUsername: "johndoe",
			setupDatabase: func(t *testing.T, db *gorm.DB) {
				t.Helper()
				// seed dữ liệu
				
				user := &model.User{
					ID:          "user-1",
					Username:    "johndoe",
					DisplayName: "John Doe",
				}
				err := db.Create(user).Error
				require.NoError(t, err)
			},
			expectedUser: &model.User{
				ID:          "user-1",
				Username:    "johndoe",
				DisplayName: "John Doe",
			},
			expectedError: false,
		},
		{
			name:          "User not found",
			inputUsername: "nonexistent",
			setupDatabase: func(t *testing.T, db *gorm.DB) {
				t.Helper()
				// không seed gì cả
			},
			expectedUser:  nil,
			expectedError: true,
		},
	}

	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := sqldb.InitMockDB(t)
			db.AutoMigrate(&model.User{})
			tc.setupDatabase(t, db)
			repo := user.NewSqlRepository(db)

			result, err := repo.GetUserByUsername(ctx, tc.inputUsername)

			// Check error
			if tc.expectedError {
				require.Error(t, err, "Expected an error")
			} else {
				require.NoError(t, err, "Unexpected error")
			}

			// Check user data if no error
			if !tc.expectedError {
				assert.Equal(t, tc.expectedUser.ID, result.ID)
				assert.Equal(t, tc.expectedUser.Username, result.Username)
				assert.Equal(t, tc.expectedUser.DisplayName, result.DisplayName)
			}
		})
	}
}
