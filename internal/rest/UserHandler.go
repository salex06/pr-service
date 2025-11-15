// Package rest - пакет с контроллерами,
// выполняющими обработку запросов и отправку ответов
package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/salex06/pr-service/internal/dto"
	"github.com/salex06/pr-service/internal/service"
)

// UserHandler представляет собой контроллер,
// который отвечает за получение запросов, связанных с сотрудниками,
// передачу на обработку в сервисы и формирование ответа
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler конструирует и возвращает объект UserHandler
func NewUserHandler(us *service.UserService) *UserHandler {
	return &UserHandler{
		userService: us,
	}
}

// HandleSetIsActiveRequest обрабатывает и формирует ответ на запрос
// изменения состояния сотрудника (активное/неактивное состояние)
func (uh *UserHandler) HandleSetIsActiveRequest(c *gin.Context) {
	var req dto.UserShort
	parseErr := c.ShouldBindBodyWithJSON(&req)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "json parsing error",
		})
		return
	}

	resp, err := uh.userService.SetIsActive(&req)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": resp,
	})
}

// HandleGetReviewRequest обрабатывает запрос и формирует ответ на получение PR`s,
// где пользователь с идентификатором user_id назначен ревьюером
func (uh *UserHandler) HandleGetReviewRequest(c *gin.Context) {
	userID := c.Query("user_id")

	resp, err := uh.userService.GetAssignedPRs(userID)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
