package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/salex06/pr-service/internal/dto"
	"github.com/salex06/pr-service/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(us *service.UserService) *UserHandler {
	return &UserHandler{
		userService: us,
	}
}

func (uh *UserHandler) HandleSetIsActiveRequest(c *gin.Context) {
	var req dto.UserShort
	c.ShouldBindBodyWithJSON(&req)

	resp, err := uh.userService.SetIsActive(&req)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": resp,
	})
}

func (uh *UserHandler) HandleGetReviewRequest(c *gin.Context) {
	userId := c.Query("user_id")

	resp, err := uh.userService.GetAssignedPRs(userId)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
