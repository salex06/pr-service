package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/salex06/pr-service/internal/dto"
	"github.com/salex06/pr-service/internal/service"
)

// PullRequestHandler представляет собой компонент,
// который отвечает за получение запросов, связанных с pull-request`ми,
// передачу на обработку в сервисы и формирование ответа
type PullRequestHandler struct {
	prService *service.PullRequestService
}

// NewPullRequestHandler конструирует и возвращает объект PullRequestHandler
func NewPullRequestHandler(prService *service.PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{
		prService: prService,
	}
}

// HandleCreateRequest отвечает за получение и формирование ответа на запрос
// открытия нового pull-request`а и назначения на него сотрудников
func (prh *PullRequestHandler) HandleCreateRequest(c *gin.Context) {
	var req dto.CreatePullRequest
	parseErr := c.ShouldBindBodyWithJSON(&req)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "json parsing error",
		})
		return
	}

	resp, err := prh.prService.CreatePullRequest(&req)
	if err != nil {
		c.JSON(err.Status, err.Error)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"pr": resp,
	})
}

// HandleMergeRequest отвечает за получение и формирование ответа на запрос
// закрытия pull-request`а и его перевода в статус MERGED
func (prh *PullRequestHandler) HandleMergeRequest(c *gin.Context) {
	var req dto.MergePullRequest
	parseErr := c.ShouldBindBodyWithJSON(&req)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "json parsing error",
		})
		return
	}

	resp, err := prh.prService.MergePullRequest(&req)
	if err != nil {
		c.JSON(err.Status, err.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pr": resp,
	})
}

// HandleReassignRequest отвечает за получение и формирование ответа на запрос
// переназначения одного сотрудника на pull-request
func (prh *PullRequestHandler) HandleReassignRequest(c *gin.Context) {
	var req dto.ReassignPullRequest
	parseErr := c.ShouldBindBodyWithJSON(&req)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "json parsing error",
		})
		return
	}

	resp, err := prh.prService.ReassignPullRequest(&req)
	if err != nil {
		c.JSON(err.Status, err.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pr":          resp.Pr,
		"replaced_by": resp.ReplacedBy,
	})
}
