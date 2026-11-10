package handler

import (
	"errors"
	"net/http"
	"time"

	urlstorage "github.com/ThienKim52/golang-dev/internal/app/repository/url_storage"
	"github.com/ThienKim52/golang-dev/internal/app/service/link"
	"github.com/ThienKim52/golang-dev/pkg/requestutils"
	"github.com/ThienKim52/golang-dev/response"
	"github.com/gin-gonic/gin"
	log "github.com/rs/zerolog/log"
)

type LinkHandler interface {
	ShortenURL(c *gin.Context)
	Redirect(c *gin.Context)
}

// LinkHandler handles link-related requests
type shortenURL struct {
	service link.LinkService
}

// NewLinkHandler creates a new link handler
func NewLinkHandler(service link.LinkService) *shortenURL {
	return &shortenURL{
		service: service,
	}
}

type ShortenInputBody struct {
	URL string `json:"url" binding:"required,url"`
	Exp *int64 `json:"exp" binding:"required,gte=300"`
}

// ShortenURL handles POST /v1/links/shorten
// @Summary Shorten a URL
// @Description Creates a short code for a given URL
// @Tags links
// @Accept json
// @Produce json
// @Param request body handler.ShortenInputBody true "Request body"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/links/shorten [post]
func (h *shortenURL) ShortenURL(c *gin.Context) {
	req, err := requestutils.BindInputFromRequest[ShortenInputBody](c)
	if err != nil {
		return
	}

	code, err := h.service.ShortenURL(c, req.URL, time.Duration(*req.Exp)*time.Second)
	if err != nil {
		log.Error().Err(err).Str("from", "handler.shortenURL.ShortenURL").Msg("Failed to shorten URL")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Shorten URL generated successfully!",
		"code":    code,
	})
}

// Redirect Forward the request to the original url

// @Tags links
// @Accept application/json
// @Produce application/json
// @Param code path string true "Shorten code"
// @Success 302
// @Router /v1/links/redirect/{code} [get]
func (s *shortenURL) Redirect(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, response.InputErrResponse)
	}

	url, err := s.service.GetLinkFromCode(c, code)
	if err != nil {
		if errors.Is(err, urlstorage.ErrCodeNotFound) {
			c.JSON(http.StatusNotFound, response.InternalErrResponse)
			return
		}
		log.Error().Err(err).Str("from", "handler.LinkHandler.Redirect").Msg("Can't get url from code")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}
	c.Redirect(http.StatusFound, url)

}
