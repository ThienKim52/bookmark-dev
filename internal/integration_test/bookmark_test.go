package integration_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ThienKim52/golang-dev/internal/api"
	"github.com/ThienKim52/golang-dev/internal/app/model"
	"github.com/ThienKim52/golang-dev/pkg/sqldb"
	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type mockJWTValidator struct {
	claims jwt.MapClaims
	err    error
}

func (m *mockJWTValidator) ValidateJWT(tokenStr string) (jwt.MapClaims, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.claims, nil
}

func setupBookmarkTestEngine(t *testing.T, db *gorm.DB, val *mockJWTValidator) api.Engine {
	t.Helper()
	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(func() { mr.Close() })
	redisClient := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { redisClient.Close() })

	return api.New(&api.EngineOpts{
		App:         gin.New(),
		Config:      &api.Config{},
		RedisClient: redisClient,
		DB:          db,
		JWTGen:      &mockJWTGenerator{},
		JWTVal:      val,
	})
}

func TestBookmarkAPI(t *testing.T) {
	t.Parallel()

	db := sqldb.InitMockDB(t)
	err := db.AutoMigrate(&model.User{}, &model.Bookmark{})
	require.NoError(t, err)

	// Seed user
	user := &model.User{
		Base: model.Base{
			ID: "deb745af-1a62-4efa-99a0-f06b274bd990",
		},
		Username:    "testuser",
		Email:       "test@example.com",
		DisplayName: "Test User",
	}
	require.NoError(t, db.Create(user).Error)

	jwtVal := &mockJWTValidator{
		claims: jwt.MapClaims{
			"sub":   "deb745af-1a62-4efa-99a0-f06b274bd990",
			"email": "test@example.com",
		},
	}
	engine := setupBookmarkTestEngine(t, db, jwtVal)

	t.Run("Create Bookmark 1 successfully via POST /v1/bookmarks", func(t *testing.T) {
		body := `{"description":"My First Bookmark","url":"https://golang.org"}`
		req := httptest.NewRequest("POST", "/v1/bookmarks", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid-token")
		rec := httptest.NewRecorder()

		engine.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "https://golang.org")
		assert.Contains(t, rec.Body.String(), "My First Bookmark")
	})

	t.Run("Create Bookmark 2 successfully for the same user (no ID collision)", func(t *testing.T) {
		body := `{"description":"My Second Bookmark","url":"https://github.com"}`
		req := httptest.NewRequest("POST", "/v1/bookmarks", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid-token")
		rec := httptest.NewRecorder()

		engine.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "https://github.com")
		assert.Contains(t, rec.Body.String(), "My Second Bookmark")
	})

	t.Run("Create Bookmark without auth -> 401", func(t *testing.T) {
		body := `{"description":"No Auth Bookmark","url":"https://example.com"}`
		req := httptest.NewRequest("POST", "/v1/bookmarks", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		engine.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("Create Bookmark with invalid URL -> 400", func(t *testing.T) {
		body := `{"description":"Invalid URL Bookmark","url":"not-a-valid-url"}`
		req := httptest.NewRequest("POST", "/v1/bookmarks", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid-token")
		rec := httptest.NewRecorder()

		engine.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Get Bookmarks via GET /v1/bookmarks", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/bookmarks?page=1&limit=10", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		rec := httptest.NewRecorder()

		engine.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "https://golang.org")
		assert.Contains(t, rec.Body.String(), "https://github.com")
	})
}
