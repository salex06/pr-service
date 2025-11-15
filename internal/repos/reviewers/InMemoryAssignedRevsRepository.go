package reviewers

import (
	"context"
	"slices"
)

type InMemoryAssignedRevsRepository struct {
	storage    map[string][]string //userId - []pullRequestIds
	storageRev map[string][]string //pullRequestId - []userIds
}

func NewInMemoryAssignedRevsRepository() *InMemoryAssignedRevsRepository {
	return &InMemoryAssignedRevsRepository{
		storage:    make(map[string][]string),
		storageRev: make(map[string][]string),
	}
}

func (repo *InMemoryAssignedRevsRepository) GetAssignedPullRequestIds(ctx context.Context, userId string) ([]string, error) {
	return repo.storage[userId], nil
}

func (repo *InMemoryAssignedRevsRepository) CreateAssignment(ctx context.Context, userId string, prId string) error {
	repo.storage[userId] = append(repo.storage[userId], prId)
	repo.storageRev[prId] = append(repo.storageRev[prId], userId)
	return nil
}

func (repo *InMemoryAssignedRevsRepository) GetAssignedReviewersIds(ctx context.Context, prId string) ([]string, error) {
	return repo.storageRev[prId], nil
}

func (repo *InMemoryAssignedRevsRepository) DeleteAssignment(ctx context.Context, userId string, prId string) error {
	repo.storage[userId] = slices.DeleteFunc(repo.storage[userId], func(currPrId string) bool { return prId == currPrId })
	repo.storageRev[prId] = slices.DeleteFunc(repo.storageRev[prId], func(currUserId string) bool { return currUserId == userId })

	return nil
}
