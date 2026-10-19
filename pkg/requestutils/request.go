package requestutils

import (
	"errors"
	"io"
	"net/http"

	"github.com/ThienKim52/golang-dev/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var InputValidator = validator.New(validator.WithRequiredStructEnabled())

func BindInputFromRequest[T any](c *gin.Context) (*T, error) {
	reqInput := new(T)

	if c.Request.Method != http.MethodGet {
		if err := c.ShouldBindJSON(reqInput); err != nil && !errors.Is(err, io.EOF) {
			c.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldError(err))
			return nil, err
		}
	}

	if err := c.ShouldBindUri(reqInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return nil, err
	}
	if err := c.ShouldBindQuery(reqInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return nil, err
	}
	if err := c.ShouldBindHeader(reqInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return nil, err
	}

	if err := InputValidator.Struct(reqInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return nil, err
	}

	return reqInput, nil
}