package reviewers

import "slices"

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

func (repo *InMemoryAssignedRevsRepository) GetAssignedPullRequestIds(userId string) []string {
	return repo.storage[userId]
}

func (repo *InMemoryAssignedRevsRepository) CreateAssignment(userId string, prId string) {
	repo.storage[userId] = append(repo.storage[userId], prId)
	repo.storageRev[prId] = append(repo.storageRev[prId], userId)
}

func (repo *InMemoryAssignedRevsRepository) GetAssignedReviewersIds(prId string) []string {
	return repo.storageRev[prId]
}

func (repo *InMemoryAssignedRevsRepository) DeleteAssignment(userId string, prId string) {
	repo.storage[userId] = slices.DeleteFunc(repo.storage[userId], func(currPrId string) bool { return prId == currPrId })
	repo.storageRev[prId] = slices.DeleteFunc(repo.storageRev[prId], func(currUserId string) bool { return currUserId == userId })
}
