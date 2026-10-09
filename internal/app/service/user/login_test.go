package user

import (
	"context"
	"errors"
	"testing"

	"github.com/ThienKim52/golang-dev/internal/app/model"
	mocks_userrepo "github.com/ThienKim52/golang-dev/internal/app/repository/user/mocks"
	"github.com/ThienKim52/golang-dev/pkg/dbutils"
	mocks_jwtutils "github.com/ThienKim52/golang-dev/pkg/jwtutils/mocks"
	mocks_utils "github.com/ThienKim52/golang-dev/pkg/utils/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	existingUser := &model.User{
		ID:       "user-123",
		Username: "johndoe",
		Password: "hashed-password",
		Email:    "john@test.com",
	}

	testCases := []struct {
		name          string
		username      string
		password      string
		setupMocks    func(repo *mocks_userrepo.Repository, hasher *mocks_utils.Hasher, jwtGen *mocks_jwtutils.JWTGenerator)
		expectedToken string
		expectedErr   error
	}{
		{
			name:     "Login successfully",
			username: "johndoe",
			password: "password123",
			setupMocks: func(repo *mocks_userrepo.Repository, hasher *mocks_utils.Hasher, jwtGen *mocks_jwtutils.JWTGenerator) {
				ctx := context.Background()
				repo.On("GetUserByUsername", ctx, "johndoe").Return(existingUser, nil).Once()
				hasher.On("Compare", "hashed-password", "password123").Return(true).Once()
				jwtGen.On("GenerateJWT", mock.MatchedBy(func(claims jwt.MapClaims) bool {
					return claims["sub"] == "user-123" && claims["email"] == "john@test.com"
				})).Return("jwt-token-abc", nil).Once()
			},
			expectedToken: "jwt-token-abc",
			expectedErr:   nil,
		},
		{
			name:     "User not found",
			username: "nonexistent",
			password: "password123",
			setupMocks: func(repo *mocks_userrepo.Repository, hasher *mocks_utils.Hasher, jwtGen *mocks_jwtutils.JWTGenerator) {
				ctx := context.Background()
				repo.On("GetUserByUsername", ctx, "nonexistent").Return(nil, dbutils.ErrRecordNotFound).Once()
			},
			expectedToken: "",
			expectedErr:   ErrInvalidCredentials,
		},
		{
			name:     "Database error",
			username: "johndoe",
			password: "password123",
			setupMocks: func(repo *mocks_userrepo.Repository, hasher *mocks_utils.Hasher, jwtGen *mocks_jwtutils.JWTGenerator) {
				ctx := context.Background()
				repo.On("GetUserByUsername", ctx, "johndoe").Return(nil, errors.New("connection refused")).Once()
			},
			expectedToken: "",
			expectedErr:   errors.New("connection refused"),
		},
		{
			name:     "Wrong password",
			username: "johndoe",
			password: "wrongpassword",
			setupMocks: func(repo *mocks_userrepo.Repository, hasher *mocks_utils.Hasher, jwtGen *mocks_jwtutils.JWTGenerator) {
				ctx := context.Background()
				repo.On("GetUserByUsername", ctx, "johndoe").Return(existingUser, nil).Once()
				hasher.On("Compare", "hashed-password", "wrongpassword").Return(false).Once()
			},
			expectedToken: "",
			expectedErr:   ErrInvalidCredentials,
		},
		{
			name:     "JWT generation fails",
			username: "johndoe",
			password: "password123",
			setupMocks: func(repo *mocks_userrepo.Repository, hasher *mocks_utils.Hasher, jwtGen *mocks_jwtutils.JWTGenerator) {
				ctx := context.Background()
				repo.On("GetUserByUsername", ctx, "johndoe").Return(existingUser, nil).Once()
				hasher.On("Compare", "hashed-password", "password123").Return(true).Once()
				jwtGen.On("GenerateJWT", mock.Anything).Return("", errors.New("jwt signing error")).Once()
			},
			expectedToken: "",
			expectedErr:   errors.New("jwt signing error"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repoMock := mocks_userrepo.NewRepository(t)
			hasherMock := mocks_utils.NewHasher(t)
			jwtGenMock := mocks_jwtutils.NewJWTGenerator(t)

			tc.setupMocks(repoMock, hasherMock, jwtGenMock)

			srv := NewService(repoMock, hasherMock, jwtGenMock)
			token, err := srv.Login(context.Background(), tc.username, tc.password)

			assert.Equal(t, tc.expectedToken, token)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}