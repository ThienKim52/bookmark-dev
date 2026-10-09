package user

import (
	"net/http"

	"github.com/ThienKim52/golang-dev/pkg/requestutils"
	"github.com/ThienKim52/golang-dev/response"
	"github.com/gin-gonic/gin"
)
type AllTypeInput struct {
	ID   string `uri:"id" validate:"required"`
	Name string `uri:"name" validate:"required"`
	Sort   string `form:"sort" validate:"required"`
	Type string `form:"type" validate:"required"`
	Count string `header:"count" validate:"required"`
}

func (h *userHandler) Test(c *gin.Context) {
	input, err := requestutils.BindInputFromRequest[AllTypeInput](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}
	println(input.ID, input.Count)
}