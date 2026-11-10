package bookmark

import (
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/ThienKim52/golang-dev/internal/app/handler/dto"
	"github.com/ThienKim52/golang-dev/pkg/dbutils"
	"github.com/ThienKim52/golang-dev/pkg/requestutils"
	"github.com/ThienKim52/golang-dev/response"
	"github.com/gin-gonic/gin"
)

// UpdateBookmark godoc
// @Summary Update a bookmark
// @Description Update a bookmark belonging to the authenticated user.
// @Tags bookmarks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Bookmark ID"
// @Param input body dto.BookmarkUpdateInput true "Bookmark update input"
// @Success 200 {object} response.Message
// @Router /v1/bookmarks/{id} [put]
func (h *bookmarkHandler) UpdateBookmark(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.InputErrResponse)
		return 
	}
	input, uid, err := requestutils.BindInputFromRequestWithAuth[dto.BookmarkUpdateInput](c)
	if err != nil {
		return
	}
	err = h.bookmarkSvc.UpdateBookmark(c, uid, id, input.Description, input.URL)
	if err != nil {
		log.Error().Err(err).Str("userID", uid).Str("bookmarkID", id).Msg("Failed to update bookmark")
		switch {
		case errors.Is(err, dbutils.ErrDuplicationEmail):
			c.AbortWithStatusJSON(http.StatusConflict, response.Message{
				Message: "Email already taken",
			})
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