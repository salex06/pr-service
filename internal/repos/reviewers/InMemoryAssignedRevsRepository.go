package reviewers

import (
	"context"
	"slices"

	"github.com/salex06/pr-service/internal/dto"
)

// InMemoryAssignedRevsRepository представляет собой компонент,
// отвечающий за взаимодействие с in-memory хранилищем (map),
// где находится информация о назначениях сотрудников на PR's
type InMemoryAssignedRevsRepository struct {
	storage    map[string][]string // userId - []pullRequestIds
	storageRev map[string][]string // pullRequestId - []userIds
}

// NewInMemoryAssignedRevsRepository конструирует и возвращает объект InMemoryAssignedRevsRepository
func NewInMemoryAssignedRevsRepository() *InMemoryAssignedRevsRepository {
	return &InMemoryAssignedRevsRepository{
		storage:    make(map[string][]string),
		storageRev: make(map[string][]string),
	}
}

// GetAssignedPullRequestIds возвращает слайс идентификаторов
// PR`s, на которые назначен сотрудник с идентификатором userID
func (repo *InMemoryAssignedRevsRepository) GetAssignedPullRequestIds(ctx context.Context, userID string) ([]string, error) {
	return repo.storage[userID], nil
}

// CreateAssignment сохраняет назначение сотрудника с
// идентификатором userID на PR с идентификатором prID
func (repo *InMemoryAssignedRevsRepository) CreateAssignment(ctx context.Context, userID, prID string) error {
	repo.storage[userID] = append(repo.storage[userID], prID)
	repo.storageRev[prID] = append(repo.storageRev[prID], userID)
	return nil
}

// GetAssignedReviewersIds возвращает слайс идентификаторов
// сотрудников, которые назначены на PR с идентификатором prID
func (repo *InMemoryAssignedRevsRepository) GetAssignedReviewersIds(ctx context.Context, prID string) ([]string, error) {
	return repo.storageRev[prID], nil
}

// DeleteAssignment удаляет назначение сотрудника
// с идентификатором userID на PR с идентификатором prID
func (repo *InMemoryAssignedRevsRepository) DeleteAssignment(ctx context.Context, userID, prID string) error {
	repo.storage[userID] = slices.DeleteFunc(repo.storage[userID], func(currPrId string) bool { return prID == currPrId })
	repo.storageRev[prID] = slices.DeleteFunc(repo.storageRev[prID], func(currUserId string) bool { return currUserId == userID })

	return nil
}

// GetAssignmentsCountByReviewerID возвращает набор пар
// "идентификатор ревьюера - количество назначений на PR данного пользователя"
func (repo *InMemoryAssignedRevsRepository) GetAssignmentsCountByReviewerID(ctx context.Context) ([]*dto.AssignmentsByUser, error) {
	assignmentsByUsers := make([]*dto.AssignmentsByUser, 0, len(repo.storage))
	for k, v := range repo.storage {
		assignmentsByUsers = append(assignmentsByUsers, &dto.AssignmentsByUser{
			UserID:           k,
			AssignmentsCount: len(v),
		})
	}

	return assignmentsByUsers, nil
}
