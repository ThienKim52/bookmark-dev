package user

import (
	"errors"
	"net/http"

	"github.com/ThienKim52/golang-dev/internal/app/service/user"
	"github.com/ThienKim52/golang-dev/pkg/requestutils"
	"github.com/ThienKim52/golang-dev/response"
	"github.com/gin-gonic/gin"
	log "github.com/rs/zerolog/log"
)

type loginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// Login handles user login requests.
// @Summary Login user
// @Description Authenticate user with username and password, returns JWT token
// @Tags User
// @Accept json
// @Produce json
// @Param body body registerInput true "User login credentials"
// @Router /v1/users/login [post]
func (h *userHandler) Login(c *gin.Context) {
	input, err := requestutils.BindInputFromRequest[loginInput](c)
	if err != nil {
		return
	}
	
	token, err := h.svc.Login(c, input.Username, input.Password)
	switch {
	case errors.Is(err, user.ErrInvalidCredentials):
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Message{Message: "Invalid credentials"})
		return
	case err == nil:
	default:
		log.Err(err).Msg("Failed to register user")
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in successfully!",
		"data":    token,
	})

}

