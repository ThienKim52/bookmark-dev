package integration_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ThienKim52/golang-dev/internal/api"
	"github.com/ThienKim52/golang-dev/internal/app/model"
	"github.com/redis/go-redis/v9"  
	"github.com/ThienKim52/golang-dev/pkg/sqldb"
	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type mockJWTGenerator struct{}

func (m *mockJWTGenerator) GenerateJWT(claims jwt.MapClaims) (string, error) {
	return "mock-jwt-token-test", nil
}

func setupLoginTestEngine(t *testing.T, db *gorm.DB) api.Engine {
	t.Helper()
	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(func() {mr.Close()})
	redisClient := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { redisClient.Close()})
	return api.New(&api.EngineOpts{
		App: gin.New(),
		Config: &api.Config{},
		RedisClient: redisClient,
		DB: db,
		JWTGen: &mockJWTGenerator{},
	})
}

func seedUserWithHashPwd(t *testing.T, db *gorm.DB, username, plainPwassword, email string) {
	t.Helper()
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(plainPwassword), bcrypt.DefaultCost)
	require.NoError(t, err)
	user := &model.User{
		ID: "test1",
		Username: username,
		Password: string(hashBytes),
		Email: email,
		DisplayName: "test",
	}
	err = db.Session(&gorm.Session{SkipHooks: true}).Create(user).Error
	require.NoError(t, err)
}

func TestLogin(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name               string
		reqBody string
		setupDB func(t *testing.T, db *gorm.DB)
		expectedStatusCode int
		expectedMessage    string
		verifyResponse     func(t *testing.T, body string) // verify response
	}{
		{
			name: "Login successfully",
			reqBody: `{"username":"johndoe","password":"password123"}`,
			setupDB: func(t *testing.T, db *gorm.DB) {
				seedUserWithHashPwd(t, db, "johndoe", "password123", "john@test.com")
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage: "Logged in successfully!",
			verifyResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, "mock-jwt-token-test")
			},
		},
		{
            name:    "Wrong password",
            reqBody: `{"username":"johndoe","password":"wrongpassword"}`,
            setupDB: func(t *testing.T, db *gorm.DB) {
                seedUserWithHashPwd(t, db, "johndoe", "password123", "john@test.com")
            },
            expectedStatusCode: http.StatusBadRequest,
            expectedMessage:    "Invalid credentials",
        },
        {
            name:    "Username not exist",
            reqBody: `{"username":"nonexistent","password":"password123"}`,
            setupDB: func(t *testing.T, db *gorm.DB) {
                // không seed gì cả
            },
            expectedStatusCode: http.StatusBadRequest,
            expectedMessage:    "Invalid credentials",
        },
        {
            name:    "Password < 8",
            reqBody: `{"username":"johndoe","password":"short"}`,
            setupDB: func(t *testing.T, db *gorm.DB) {},
            expectedStatusCode: http.StatusBadRequest,
        },
        {
            name:    "Missing username",
            reqBody: `{"password":"password123"}`,
            setupDB: func(t *testing.T, db *gorm.DB) {},
            expectedStatusCode: http.StatusBadRequest,
        },
        {
            name:    "Body JSON invalid",
            reqBody: `not-a-json`,
            setupDB: func(t *testing.T, db *gorm.DB) {},
            expectedStatusCode: http.StatusBadRequest,
        },
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := sqldb.InitMockDB(t)
			err := db.AutoMigrate(&model.User{})
			require.NoError(t, err)

			if tc.setupDB != nil {
				tc.setupDB(t, db)
			}

			apiEngine := setupLoginTestEngine(t, db)
			
			req := httptest.NewRequest("POST", "/v1/users/login", bytes.NewBufferString(tc.reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			apiEngine.ServeHTTP(rec, req)
			if tc.expectedMessage != "" {
                assert.Contains(t, rec.Body.String(), tc.expectedMessage)
            }
            if tc.verifyResponse != nil {
                tc.verifyResponse(t, rec.Body.String())
            }
		})
	}
}
