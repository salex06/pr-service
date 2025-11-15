package pr

import (
	"context"

	"github.com/salex06/pr-service/internal/model"
)

type PullRequestRepository interface {
	PullRequestExists(ctx context.Context, prId string) (bool, error)

	GetPullRequest(ctx context.Context, prId string) (*model.PullRequest, error)
	GetPullRequests(ctx context.Context, prIds []string) ([]*model.PullRequest, error)

	SavePullRequest(ctx context.Context, pr *model.PullRequest) error
	UpdatePullRequest(ctx context.Context, pr *model.PullRequest) error
}
