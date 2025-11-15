package dto

import "github.com/salex06/pr-service/internal/entity"

// PullRequestShort является формой представления сущности PullRequest
// с уникальным идентификатором, именем, именем автора и статусом
type PullRequestShort struct {
	PullRequestID   string                   `json:"pull_request_id"`
	PullRequestName string                   `json:"pull_request_name"`
	AuthorID        string                   `json:"author_id"`
	Status          entity.PullRequestStatus `json:"status"`
}
