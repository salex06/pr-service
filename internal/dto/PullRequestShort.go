package dto

import "github.com/salex06/pr-service/internal/model"

type PullRequestShort struct {
	PullRequestId   string                  `json:"pull_request_id"`
	PullRequestName string                  `json:"pull_request_name"`
	AuthorId        string                  `json:"author_id"`
	Status          model.PullRequestStatus `json:"status"`
}
