package pr

import "github.com/salex06/pr-service/internal/model"

type InMemoryPullRequestRepository struct {
	storage map[string]*model.PullRequest
}

func NewInMemoryPullRequestRepository() *InMemoryPullRequestRepository {
	return &InMemoryPullRequestRepository{
		storage: make(map[string]*model.PullRequest),
	}
}

func (repo *InMemoryPullRequestRepository) GetPullRequest(prId string) *model.PullRequest {
	return repo.storage[prId]
}

func (repo *InMemoryPullRequestRepository) GetPullRequests(prIds []string) []*model.PullRequest {
	prs := make([]*model.PullRequest, 0, len(prIds))
	for _, v := range prIds {
		prs = append(prs, repo.GetPullRequest(v))
	}
	return prs
}

func (repo *InMemoryPullRequestRepository) PullRequestExists(prId string) bool {
	_, ok := repo.storage[prId]
	return ok
}

func (repo *InMemoryPullRequestRepository) SavePullRequest(pr *model.PullRequest) *model.PullRequest {
	repo.storage[pr.PullRequestId] = pr
	return repo.storage[pr.PullRequestId]
}

func (repo *InMemoryPullRequestRepository) UpdatePullRequest(pr *model.PullRequest) *model.PullRequest {
	repo.storage[pr.PullRequestId] = pr
	return repo.storage[pr.PullRequestId]
}
