// Package reviewers - пакет с репозиториями, отвечающими за взаимодействие с БД,
// где хранится информация о назначенных на PR сотрудниках
package reviewers

import (
	"context"

	"github.com/salex06/pr-service/internal/dto"
)

// AssignedRevsRepository представляет собой интерфейс взаимодействия
// с различными базами данных и хранилищами, где содержится информация
// о назначениях сотрудников на PR's
type AssignedRevsRepository interface {
	GetAssignedPullRequestIds(ctx context.Context, userID string) ([]string, error)
	GetAssignedReviewersIds(ctx context.Context, pullRequestID string) ([]string, error)

	GetAssignmentsCountByReviewerID(context.Context) ([]*dto.AssignmentsByUser, error)

	CreateAssignment(ctx context.Context, userID string, prID string) error
	DeleteAssignment(ctx context.Context, userID string, prID string) error
}
