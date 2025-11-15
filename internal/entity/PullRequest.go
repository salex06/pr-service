package entity

import "time"

// PullRequestStatus представляет тип,
// определяющий статус PR - открыт(OPEN)/закрыт(MERGED)
type PullRequestStatus string

// Константы, определяющие допустимые типы статуса PR (открыт/закрыт)
const (
	OPEN   PullRequestStatus = "OPEN"
	MERGED PullRequestStatus = "MERGED"
)

// PullRequest представляет сущность
// с идентификатором, названием, автором, статусом и
// набором назначенных ревьюеров (до 2)
type PullRequest struct {
	PullRequestID   string
	PullRequestName string
	AuthorID        string
	Status          PullRequestStatus
	CreatedAt       *time.Time
	MergedAt        *time.Time
}
