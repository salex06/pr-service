package dto

// MergePullRequest определяет структуру запроса
// на закрытие PR и перевода его в статус MERGED
type MergePullRequest struct {
	PullRequestID string `json:"pull_request_id"`
}
