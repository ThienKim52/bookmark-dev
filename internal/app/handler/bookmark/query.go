package bookmark

import (
	"net/http"

	"github.com/ThienKim52/golang-dev/internal/app/handler/dto"
	"github.com/ThienKim52/golang-dev/internal/app/model"
	"github.com/ThienKim52/golang-dev/pkg/requestutils"
	"github.com/ThienKim52/golang-dev/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// GetBookmarks godoc
// @Summary List bookmarks
// @Description Get a paginated list of bookmarks belonging to the authenticated user.
// @Tags bookmarks
// @Security BearerAuth
// @Produce json
// @Param page query integer false "Page number (minimum 1)" minimum(1) default(1)
// @Param limit query integer false "Number of items per page" minimum(1) default(20)
// @Success 200 {object} dto.SuccessResponse[[]model.Bookmark] "Bookmarks retrieved successfully"
// @Failure 400 {object} object "Invalid pagination parameters"
// @Failure 401 {object} object "Unauthorized"
// @Failure 500 {object} object "Internal server error"
// @Router /v1/bookmarks [get]
func (h *bookmarkHandler) GetBookmarks(c *gin.Context) {
	input, uid, err := requestutils.BindInputFromRequestWithAuth[dto.BookmarkPaginationInput](c)
	if err != nil {
		return
	}
	if input.Page == 0 {
		input.Page = 1
	}
	if input.Limit == 0 {
		input.Limit = 20
	}

	if input.Page < 1 || input.Limit < 1 {
		c.JSON(http.StatusBadRequest, response.InputErrResponse)
		return
	}

	res, err := h.bookmarkSvc.GetBookmarks(c, uid, input.Page, input.Limit)
	if err != nil {
		log.Error().
			Err(err).
			Str("userID", uid).
			Msg("GetBookmarks error")

		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}
	c.JSON(http.StatusOK, &dto.SuccessResponse[[]*model.Bookmark]{
		Data: res.Bookmarks,
		Pagination: &dto.Pagination{
			Page:  input.Page,
			Limit: input.Limit,
			Total: res.Total,
		},
	})

}
