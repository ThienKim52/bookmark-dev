package requestutils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var ErrNoClaims = errors.New("No claims found in request")
func GetUserIDFromRequest(c *gin.Context) (string, error) {
	claims, ok := c.Get("claims")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return "", ErrNoClaims
	}
	tokenInfo, ok := claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return "", ErrNoClaims
	}
	uid, ok := tokenInfo["sub"].(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return "", ErrNoClaims
	}
	return uid, nil
}