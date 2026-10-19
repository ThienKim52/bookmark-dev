package user

import (
	"net/http"
	"time"

	"errors"

	"github.com/ThienKim52/golang-dev/pkg/dbutils"
	"github.com/ThienKim52/golang-dev/pkg/requestutils"
	"github.com/ThienKim52/golang-dev/response"
	"github.com/gin-gonic/gin"
	log "github.com/rs/zerolog/log"
)

type registerInput struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required,min=8"`
	DisplayName string `json:"display_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
}

type userResponseData struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type registerResponse struct {
	Data    userResponseData `json:"data"`
	Message string           `json:"message"`
}
// Register handles user registration requests.
// It validates the JSON input, delegates to the service layer for user creation,
// and returns the created user or an appropriate error response.

// @Summary Register a new user
// @Description Register a new user with the provided information
// @Tags User
// @Accept json
// @Produce json
// @Param body body registerInput true "User registration details"
// @Router /v1/users/register [post]
func (h *userHandler) Register(c *gin.Context) {
	input, err := requestutils.BindInputFromRequest[registerInput](c)
	if err != nil {
		return
	}
	createdUser, err := h.svc.CreateUser(c, input.Username, input.Password, input.DisplayName, input.Email)
	if err != nil {
		log.Error().Err(err).Msg("Failed to register user")
		switch{
		case errors.Is(err, dbutils.ErrDuplicationUsername):
		c.AbortWithStatusJSON(http.StatusConflict, response.Message{
		Message: "Username already taken",
		})
		return
		case errors.Is(err, dbutils.ErrDuplicationEmail):
		c.AbortWithStatusJSON(http.StatusConflict, response.Message{
		Message: "Email already taken",
		})
		return

		case err == nil:
		default:
		log.Err(err).Msg("Failed to register user")
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
		}
	}
	res := registerResponse{
		Data: userResponseData{
			ID:          createdUser.ID,
			Username:    createdUser.Username,
			Email:       createdUser.Email,
			DisplayName: createdUser.DisplayName,
			CreatedAt:   createdUser.CreatedAt,
			UpdatedAt:   createdUser.UpdatedAt,
		},
		Message: "Register an user successfully!",
	}
	c.JSON(http.StatusOK, res)
}
