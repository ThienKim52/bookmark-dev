package user

import (
	"net/http"

	"github.com/ThienKim52/golang-dev/pkg/requestutils"
	"github.com/gin-gonic/gin"
)

// GetSelfInfo get your current information
// @Summary get your current information
// @Description get your current information
// @Tags user
// @Security BearerAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} object{data=model.User} "Success"
// @Router /v1/self/info [get]
func (h *userHandler) GetUserInfo(c *gin.Context) {
	uid, err := requestutils.GetUserIDFromRequest(c)
	if err != nil {
		return 
	}

	curUser, err := h.svc.GetUserInfo(c, uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
        "data": gin.H{
            "id":          curUser.ID,
            "username":    curUser.Username,
            "email":       curUser.Email,
            "displayName": curUser.DisplayName,
			"created_at":  curUser.CreatedAt,
        	"updated_at":  curUser.UpdatedAt,
        },
    })
	
}