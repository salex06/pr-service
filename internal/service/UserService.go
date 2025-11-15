// Package service - пакет с компонентами,
// отвечающими за выполнение бизнес-логики приложения
package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/salex06/pr-service/internal/converter"
	"github.com/salex06/pr-service/internal/dto"
	prRepos "github.com/salex06/pr-service/internal/repos/pr"
	revsRepos "github.com/salex06/pr-service/internal/repos/reviewers"
	userRepos "github.com/salex06/pr-service/internal/repos/user"
)

// UserService представляет компонент,
// отвечающий за выполнение бизнес-логики приложения,
// связанной с пользователями (изменение состояния, получение связанных PR)
type UserService struct {
	userRepository         *userRepos.UserRepository
	assignedRevsRepository *revsRepos.AssignedRevsRepository
	pullRequestRepository  *prRepos.PullRequestRepository
}

// NewUserService конструирует и возвращает объект структуры UserService
func NewUserService(
	ur *userRepos.UserRepository,
	ar *revsRepos.AssignedRevsRepository,
	pr *prRepos.PullRequestRepository) *UserService {
	return &UserService{
		userRepository:         ur,
		assignedRevsRepository: ar,
		pullRequestRepository:  pr,
	}
}

// SetIsActive изменяет состояние сотрудника (активен или нет)
func (us *UserService) SetIsActive(req *dto.UserShort) (*dto.User, *dto.ErrorResponse) {
	if user, _ := (*us.userRepository).GetUser(context.Background(), req.UserID); user != nil {
		user.IsActive = req.IsActive
		err := (*us.userRepository).UpdateUser(context.Background(), user)
		if err != nil {
			return nil, &dto.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error: map[string]string{
					"code":    "INTERNAL_ERROR",
					"message": fmt.Sprintf("unable update user: %s", err),
				},
			}
		}

		return converter.ConvertUserEntityToDto(user), nil
	}

	return nil, &dto.ErrorResponse{
		Status: http.StatusNotFound,
		Error: map[string]string{
			"code":    string(dto.NotFound),
			"message": "resource not found",
		},
	}
}

// GetAssignedPRs возвращает пулл-реквесты,
// на которые назначен сотрудник с идентификатором userID
func (us *UserService) GetAssignedPRs(userID string) (*dto.AssignedPullRequests, *dto.ErrorResponse) {
	if exists, _ := (*us.userRepository).UserExists(context.Background(), userID); exists {
		prIds, _ := (*us.assignedRevsRepository).GetAssignedPullRequestIds(context.Background(), userID)
		prs, _ := (*us.pullRequestRepository).GetPullRequests(context.Background(), prIds)

		return converter.ConvertPRsToAssignedPRs(userID, prs), nil
	}

	// В API не прописана данная ветка
	return nil, &dto.ErrorResponse{
		Status: http.StatusNotFound,
		Error: map[string]string{
			"code":    string(dto.NotFound),
			"message": "resource not found",
		},
	}
}
