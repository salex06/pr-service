package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/salex06/pr-service/internal/dto"
	"github.com/salex06/pr-service/internal/service"
)

type PullRequestHandler struct {
	prService *service.PullRequestService
}

func NewPullRequestHandler(prService *service.PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{
		prService: prService,
	}
}

func (prh *PullRequestHandler) HandleCreateRequest(c *gin.Context) {
	var req dto.CreatePullRequest
	c.ShouldBindBodyWithJSON(&req)

	resp, err := prh.prService.CreatePullRequest(&req)
	if err != nil {
		c.JSON(err.Status, err.Error)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"pr": resp,
	})
}

func (prh *PullRequestHandler) HandleMergeRequest(c *gin.Context) {
	var req dto.MergePullRequest
	c.ShouldBindBodyWithJSON(&req)

	resp, err := prh.prService.MergePullRequest(&req)
	if err != nil {
		c.JSON(err.Status, err.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pr": resp,
	})
}

func (prh *PullRequestHandler) HandleReassignRequest(c *gin.Context) {
	var req dto.ReassignPullRequest
	c.ShouldBindBodyWithJSON(&req)

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
