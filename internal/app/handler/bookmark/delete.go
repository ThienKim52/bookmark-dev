package bookmark

import (
	"errors"
	"net/http"
	"github.com/rs/zerolog/log"
	"github.com/ThienKim52/golang-dev/pkg/dbutils"
	"github.com/ThienKim52/golang-dev/pkg/requestutils"
	"github.com/ThienKim52/golang-dev/response"
	"github.com/gin-gonic/gin"
)

// DeleteBookmark godoc
// @Summary Delete a bookmark
// @Description Delete a bookmark belonging to the authenticated user.
// @Tags bookmarks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Bookmark ID"
// @Success 200 {object} response.Message
// @Router /v1/bookmarks/{id} [delete]
func (h *bookmarkHandler) DeleteBookmark(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.InputErrResponse)
		return 
	}
	uid, err := requestutils.GetUIDFromRequest(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Message{
			Message: "Unauthorized",
		})
		return
	}
	err = h.bookmarkSvc.DeleteBookmark(c, uid, id)
	if err != nil {
		log.Error().Err(err).Str("userID", uid).Str("bookmarkID", id).Msg("Failed to delete bookmark")
		switch {
		case errors.Is(err, dbutils.ErrBookmarkNotFound):
			c.AbortWithStatusJSON(http.StatusNotFound, response.Message{
				Message: "Bookmark not found",
			})
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		}
		return
	}
	c.JSON(http.StatusOK, response.Message{Message: "Success"})
}