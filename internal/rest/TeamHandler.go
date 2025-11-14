package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/salex06/pr-service/internal/dto"
	"github.com/salex06/pr-service/internal/service"
)

type TeamHandler struct {
	teamService *service.TeamService
}

func NewTeamHandler(ts *service.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: ts,
	}
}

func (th *TeamHandler) HandleAddTeamRequest(c *gin.Context) {
	var req dto.Team
	c.ShouldBindBodyWithJSON(&req)

	resp, err := th.teamService.AddTeam(&req)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"team": resp,
	})
}

func (th *TeamHandler) HandleGetTeamRequest(c *gin.Context) {
	teamId := c.Query("team_name")

	resp, err := th.teamService.GetTeam(teamId)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
