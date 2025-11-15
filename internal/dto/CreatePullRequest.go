package dto

// CreatePullRequest представляет структуру запроса
// на создание PR с уникальным идентификатором,
// именем и идентификатором автора
type CreatePullRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}
