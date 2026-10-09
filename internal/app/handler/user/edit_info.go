package user

import (
	"errors"
	"net/http"

	"github.com/ThienKim52/golang-dev/pkg/dbutils"
	"github.com/ThienKim52/golang-dev/pkg/requestutils"
	"github.com/ThienKim52/golang-dev/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type editSelfInfoInput struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email" binding:"email"`
}

func (h *userHandler) EditSelfInfo(c *gin.Context) {
	uid, err := requestutils.GetUserIDFromRequest(c)
	if err != nil {
		return 
	}

	input, err := requestutils.BindInputFromRequest[editSelfInfoInput](c)
	if err != nil {
		return
	}

	err = h.svc.EditCurrentUser(c, uid, input.DisplayName, input.Email)
	if err != nil {
		log.Error().Err(err).Msg("Failed to edit user info")
		switch {
		case errors.Is(err, dbutils.ErrDuplicationEmail):
			c.AbortWithStatusJSON(http.StatusConflict, response.Message{Message: "Email already taken"})
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		}
		return
	}
	c.JSON(http.StatusOK, response.Message{Message: "User info updated successfully"})
}