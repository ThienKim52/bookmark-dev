package user

import (
	"github.com/ThienKim52/golang-dev/internal/app/service/user"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetUserInfo(c *gin.Context) 
	EditSelfInfo(c *gin.Context)
	Test(c *gin.Context)
}

type userHandler struct { 
	svc user.Service
}

func NewHandler(svc user.Service) Handler {
	return &userHandler{svc: svc}
}