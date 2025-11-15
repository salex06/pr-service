package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/salex06/pr-service/internal/dto"
	"github.com/salex06/pr-service/internal/service"
)

// TeamHandler представляет контроллер,
// который отвечает за получение запросов, связанных с командами,
// передачу на обработку в сервисы и формирование ответа
type TeamHandler struct {
	teamService *service.TeamService
}

// NewTeamHandler конструирует и возвращает объект TeamHandler
func NewTeamHandler(ts *service.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: ts,
	}
}

// HandleAddTeamRequest ответчает за получение и формирование ответа
// на запрос добавления новой команды и создания/обновления сотрудников
func (th *TeamHandler) HandleAddTeamRequest(c *gin.Context) {
	var req dto.Team
	parseErr := c.ShouldBindBodyWithJSON(&req)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "json parsing error",
		})
		return
	}

	resp, err := th.teamService.AddTeam(&req)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"team": resp,
	})
}

// HandleGetTeamRequest ответчает за получение и формирование ответа
// на запрос получения информации о команде с именем team_name
func (th *TeamHandler) HandleGetTeamRequest(c *gin.Context) {
	teamID := c.Query("team_name")

	resp, err := th.teamService.GetTeam(teamID)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
