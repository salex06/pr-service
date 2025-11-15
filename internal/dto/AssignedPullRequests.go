package dto

// AssignedPullRequests представляет структуру
// с уникальным идентификатором сотрудника и
// краткой информацией о PR's, на которые он назначен
type AssignedPullRequests struct {
	UserID       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}
