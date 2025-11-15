package dto

import (
	"time"

	"github.com/salex06/pr-service/internal/entity"
)

// MaxAssignedReviewers - максимальное количество ревьюеров,
// которое может быть назначено на PR
const MaxAssignedReviewers = 2

// PullRequest является формой представления сущности PullRequest
// с идентификатором, названием, идентификатором автора, статусом,
// назначенными сотрудниками, временем создания PR и временем его закрытия
type PullRequest struct {
	PullRequestID     string                   `json:"pull_request_id"`
	PullRequestName   string                   `json:"pull_request_name"`
	AuthorID          string                   `json:"author_id"`
	Status            entity.PullRequestStatus `json:"status"`
	AssignedReviewers []string                 `json:"assigned_reviewers"`
	CreatedAt         *time.Time               `json:"createdAt,omitempty"`
	MergedAt          *time.Time               `json:"mergedAt,omitempty"`
}
