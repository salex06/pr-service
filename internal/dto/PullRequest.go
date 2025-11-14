package dto

import (
	"time"

	"github.com/salex06/pr-service/internal/model"
)

const MAX_ASSIGNED_REVIEWERS = 2

type PullRequest struct {
	PullRequestId     string                  `json:"pull_request_id"`
	PullRequestName   string                  `json:"pull_request_name"`
	AuthorId          string                  `json:"author_id"`
	Status            model.PullRequestStatus `json:"status"`
	AssignedReviewers []string                `json:"assigned_reviewers"`
	CreatedAt         *time.Time              `json:"createdAt,omitempty"`
	MergedAt          *time.Time              `json:"mergedAt,omitempty"`
}
