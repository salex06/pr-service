// Package pr - пакет с репозиториями, отвечающими за взаимодействие с БД,
// где хранится информация о PR's
package pr

import (
	"context"

	"github.com/salex06/pr-service/internal/entity"
)

// PullRequestRepository представляет собой интерфейс взаимодействия
// с базами данных и хранилищами, где содержится информаци о PR's
type PullRequestRepository interface {
	PullRequestExists(ctx context.Context, prID string) (bool, error)

	GetPullRequest(ctx context.Context, prID string) (*entity.PullRequest, error)
	GetPullRequests(ctx context.Context, prIds []string) ([]*entity.PullRequest, error)

	SavePullRequest(ctx context.Context, pr *entity.PullRequest) error
	UpdatePullRequest(ctx context.Context, pr *entity.PullRequest) error

	GetOpenedPullRequestCount(ctx context.Context) (int, error)
	GetMergedPullRequestCount(ctx context.Context) (int, error)
}
