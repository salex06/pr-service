package dto

// ReassignPullRequest определяет структуру запроса
// на переназначение сотрудника на PR
type ReassignPullRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldReviewerID string `json:"old_reviewer_id"`
}
