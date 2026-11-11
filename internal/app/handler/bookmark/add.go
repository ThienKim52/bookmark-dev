package bookmark

import (
	"errors"
	"net/http"

	"github.com/ThienKim52/golang-dev/internal/app/handler/dto"
	"github.com/ThienKim52/golang-dev/internal/app/model"
	"github.com/ThienKim52/golang-dev/pkg/dbutils"
	"github.com/ThienKim52/golang-dev/pkg/requestutils"
	"github.com/ThienKim52/golang-dev/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type BookmarkInput struct {
	Description string `json:"description" validate:"lte=255"`
	URL         string `json:"url" validate:"required,url,lte=2048"`
}

// AddBookmark create a new bookmark
// @Summary create a new bookmark
// @Description create a new bookmark
// @Tags bookmarks
// @Security BearerAuth
// @Accept application/json
// @Param input body BookmarkInput true "Input required"
// @Produce application/json
// @Success 200 {object} dto.SuccessResponse[model.Bookmark] "Success"
// @Router /v1/bookmarks [post]
func (h *bookmarkHandler) AddBookmark(c *gin.Context) {
	input, uid, err := requestutils.BindInputFromRequestWithAuth[BookmarkInput](c)
	if err != nil {
		return 
	}
	res, err := h.bookmarkSvc.AddBookmark(c, input.Description, input.URL, uid)
	if err != nil {
		switch {
		case errors.Is(err, dbutils.ErrUniqueConstraint): 
			c.JSON(http.StatusBadRequest, &response.Message{Message: err.Error()})
			return
		default:
			log.Error().Err(err).Str("userID", uid).Msg("AddBookmark Error")
			c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
			return
		}
	}
	c.JSON(http.StatusOK, &dto.SuccessResponse[*model.Bookmark]{
		Data: res,
		Message: "Create a bookmark successfully!",
	})
}