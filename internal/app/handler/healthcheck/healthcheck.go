package handler

import (
	"net/http"

	"github.com/ThienKim52/golang-dev/internal/app/service/healthcheck"
	"github.com/ThienKim52/golang-dev/response"
	"github.com/gin-gonic/gin"
	log "github.com/rs/zerolog/log"
)

type HealthCheckHandler interface{
	GetResponse(c *gin.Context)
}
// HealthCheckHandler handles health check requests
type healthCheckHandler struct {
	service healthcheck.HealthCheckService
}

// NewHealthCheck creates a new health check handler
func NewHealthCheck(service healthcheck.HealthCheckService) HealthCheckHandler {
	return &healthCheckHandler{
		service: service,
	}
}

// GetResponse handles GET /health-check
// @Summary Health check endpoint
// @Description Returns the health status of the service
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /health-check [get]
func (h *healthCheckHandler) GetResponse(c *gin.Context) {

	result, err := h.service.GetHealthCheck(c)
	if err != nil {
		log.Error().Err(err).Str("from", "handler.healthCheckHandler.GetResponse").Msg("Failed to get response")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, result)
}
