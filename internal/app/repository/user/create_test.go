package user

import (
	"testing"

	"github.com/ThienKim52/golang-dev/internal/app/model"
	"github.com/ThienKim52/golang-dev/pkg/sqldb"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSQLRepository_CreateUser(t *testing.T) {
	testCases := []struct {
		name          string
		setupDB       func(t *testing.T) *gorm.DB
		inputUserName string
		inputEmail    string
		expectedErr   error
		verifyFunc    func(db *gorm.DB, username, email string)
	}{
		{
			name: "normal case",
			setupDB: func(t *testing.T) *gorm.DB {
				db := sqldb.InitMockDB(t)
				err := db.AutoMigrate(&model.User{})
				assert.NoError(t, err)
				return db
			},
			inputUserName: "test1",
			inputEmail:    "test1@gmail.com",
			expectedErr:   nil,
			verifyFunc: func(db *gorm.DB, username, email string) {
				user := model.User{}
				err := db.First(&user, "username = ?", username).Error
				assert.NoError(t, err)
				assert.Equal(t, username, user.Username)
				assert.Equal(t, email, user.Email)
				assert.NotEmpty(t, user.ID)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			db := tc.setupDB(t)
			repo := NewSqlRepository(db)
			user, err := repo.CreateUser(ctx, &model.User{
				Username: tc.inputUserName,
				Email:    tc.inputEmail,
			})
			if tc.expectedErr == nil {
				assert.NoError(t, err)
				assert.NotNil(t, user)
			} else {
				assert.ErrorIs(t, err, tc.expectedErr)
			}
			tc.verifyFunc(db, tc.inputUserName, tc.inputEmail)
		})
	}
}