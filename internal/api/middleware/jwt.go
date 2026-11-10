package middleware

import (
	"net/http"
	"strings"

	"github.com/ThienKim52/golang-dev/pkg/jwtutils"
	"github.com/gin-gonic/gin"
)

type JWTAuth interface {
	JWTAuth() gin.HandlerFunc
}

type jwtAuth struct {
	jwtVal jwtutils.JWTValidator
}

func NewJWTAuth(jwtVal jwtutils.JWTValidator) JWTAuth {
	return &jwtAuth{
		jwtVal: jwtVal,
	}
}

func (j *jwtAuth) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer"{
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		tokenStr := parts[1]
		tokenClaims, err := j.jwtVal.ValidateJWT(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Set("claims", tokenClaims)
		c.Next()
	}
}