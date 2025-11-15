package pr

import (
	"context"

	"github.com/salex06/pr-service/internal/entity"
)

// InMemoryPullRequestRepository представляет собой компонент,
// отвечающий за взаимодействие с in-memory хранилищем (map),
// где содержится информация о PR's
type InMemoryPullRequestRepository struct {
	storage map[string]*entity.PullRequest
}

// NewInMemoryPullRequestRepository конструирует и возвращает объект InMemoryPullRequestRepository
func NewInMemoryPullRequestRepository() *InMemoryPullRequestRepository {
	return &InMemoryPullRequestRepository{
		storage: make(map[string]*entity.PullRequest),
	}
}

// GetPullRequest возвращает PR с заданным идентификатором (nil - если не найден)
func (repo *InMemoryPullRequestRepository) GetPullRequest(ctx context.Context, prID string) (*entity.PullRequest, error) {
	return repo.storage[prID], nil
}

// GetPullRequests возвращает набор PR's по заданному набору идентификаторов
func (repo *InMemoryPullRequestRepository) GetPullRequests(ctx context.Context, prIds []string) ([]*entity.PullRequest, error) {
	prs := make([]*entity.PullRequest, 0, len(prIds))
	for _, v := range prIds {
		if pr, err := repo.GetPullRequest(ctx, v); err == nil {
			prs = append(prs, pr)
		}
	}
	return prs, nil
}

// PullRequestExists выполняет проверку наличия PR
// с заданным идентификатором в хранилище и возвращает результат
func (repo *InMemoryPullRequestRepository) PullRequestExists(ctx context.Context, prID string) (bool, error) {
	_, ok := repo.storage[prID]
	return ok, nil
}

// SavePullRequest выполняет сохранение PR в хранилище
func (repo *InMemoryPullRequestRepository) SavePullRequest(ctx context.Context, pr *entity.PullRequest) error {
	repo.storage[pr.PullRequestID] = pr
	return nil
}

// UpdatePullRequest выполняет обновление PR
// (для данной реализации идентично SavePullRequest)
func (repo *InMemoryPullRequestRepository) UpdatePullRequest(ctx context.Context, pr *entity.PullRequest) error {
	repo.storage[pr.PullRequestID] = pr
	return nil
}

// GetOpenedPullRequestCount возвращает число PR в статусе OPEN
func (repo *InMemoryPullRequestRepository) GetOpenedPullRequestCount(ctx context.Context) (int, error) {
	count := 0
	for _, v := range repo.storage {
		if v.Status == entity.OPEN {
			count++
		}
	}
	return count, nil
}

// GetMergedPullRequestCount возвращает число PR в статусе MERGED
func (repo *InMemoryPullRequestRepository) GetMergedPullRequestCount(ctx context.Context) (int, error) {
	count := 0
	for _, v := range repo.storage {
		if v.Status == entity.MERGED {
			count++
		}
	}
	return count, nil
}
