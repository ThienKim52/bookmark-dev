package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ThienKim52/golang-dev/internal/app/service/genpass"
	"github.com/ThienKim52/golang-dev/response"
	log "github.com/rs/zerolog/log"
)

const passwordLength = 12

// genpass interface
type GenPass interface {
	GeneratePassword(c *gin.Context)
}

type genPass struct {
	svc genpass.GenPass
}

// Constructor
func NewGenPass(svc genpass.GenPass) GenPass {
	return &genPass{svc: svc}
}

// GeneratePassword generates a random password of the given length
func (g *genPass) GeneratePassword(c *gin.Context) {
	// call service
	pass, err := g.svc.GeneratePassword(passwordLength)
	if err != nil {
		log.Error().Err(err).Str("from", "handler.genPass.GeneratePassword").Msg("Failed to generate password")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	// return response
	c.JSON(http.StatusOK, gin.H{"password": pass})
}
