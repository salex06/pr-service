package pr

import (
	"context"

	"github.com/salex06/pr-service/internal/model"
)

type InMemoryPullRequestRepository struct {
	storage map[string]*model.PullRequest
}

func NewInMemoryPullRequestRepository() *InMemoryPullRequestRepository {
	return &InMemoryPullRequestRepository{
		storage: make(map[string]*model.PullRequest),
	}
}

func (repo *InMemoryPullRequestRepository) GetPullRequest(ctx context.Context, prId string) (*model.PullRequest, error) {
	return repo.storage[prId], nil
}

func (repo *InMemoryPullRequestRepository) GetPullRequests(ctx context.Context, prIds []string) ([]*model.PullRequest, error) {
	prs := make([]*model.PullRequest, 0, len(prIds))
	for _, v := range prIds {
		if pr, err := repo.GetPullRequest(ctx, v); err == nil {
			prs = append(prs, pr)
		}
	}
	return prs, nil
}

func (repo *InMemoryPullRequestRepository) PullRequestExists(ctx context.Context, prId string) (bool, error) {
	_, ok := repo.storage[prId]
	return ok, nil
}

func (repo *InMemoryPullRequestRepository) SavePullRequest(ctx context.Context, pr *model.PullRequest) error {
	repo.storage[pr.PullRequestId] = pr
	return nil
}

func (repo *InMemoryPullRequestRepository) UpdatePullRequest(ctx context.Context, pr *model.PullRequest) error {
	repo.storage[pr.PullRequestId] = pr
	return nil
}
