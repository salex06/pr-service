package service

import (
	"context"
	"net/http"

	"github.com/salex06/pr-service/internal/dto"
	"github.com/salex06/pr-service/internal/model"
	prRepos "github.com/salex06/pr-service/internal/repos/pr"
	revsRepos "github.com/salex06/pr-service/internal/repos/reviewers"
	userRepos "github.com/salex06/pr-service/internal/repos/user"
)

type UserService struct {
	userRepository         *userRepos.UserRepository
	assignedRevsRepository *revsRepos.AssignedRevsRepository
	pullRequestRepository  *prRepos.PullRequestRepository
}

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

func (us *UserService) SetIsActive(req *dto.UserShort) (*dto.User, *dto.ErrorResponse) {
	if user, _ := (*us.userRepository).GetUser(context.Background(), req.UserId); user != nil {
		user.IsActive = req.IsActive
		(*us.userRepository).UpdateUser(context.Background(), user)

		//TODO: лучше сделать конвертер
		return (*dto.User)(user), nil
	}

	return nil, &dto.ErrorResponse{
		Status: http.StatusNotFound,
		Error: map[string]string{
			"code":    string(dto.NOT_FOUND),
			"message": "resource not found",
		},
	}
}

func (us *UserService) GetAssignedPRs(userId string) (*dto.AssignedPullRequests, *dto.ErrorResponse) {
	if exists, _ := (*us.userRepository).UserExists(context.Background(), userId); exists {
		prIds, _ := (*us.assignedRevsRepository).GetAssignedPullRequestIds(context.Background(), userId)
		prs, _ := (*us.pullRequestRepository).GetPullRequests(context.Background(), prIds)

		return us.convertToAssignedPRs(userId, prs), nil
	}

	//В API не прописана данная ветка
	return nil, &dto.ErrorResponse{
		Status: http.StatusNotFound,
		Error: map[string]string{
			"code":    string(dto.NOT_FOUND),
			"message": "resource not found",
		},
	}
}

func (us *UserService) convertToAssignedPRs(userId string, prs []*model.PullRequest) *dto.AssignedPullRequests {
	pullRequestsShort := us.convertToShortPullRequests(prs)

	return &dto.AssignedPullRequests{
		UserId:       userId,
		PullRequests: pullRequestsShort,
	}
}

func (us *UserService) convertToShortPullRequests(prs []*model.PullRequest) []dto.PullRequestShort {
	converted := make([]dto.PullRequestShort, 0, len(prs))

	for _, v := range prs {
		converted = append(converted, dto.PullRequestShort{
			PullRequestId:   v.PullRequestId,
			PullRequestName: v.PullRequestName,
			AuthorId:        v.AuthorId,
			Status:          v.Status,
		})
	}

	return converted
}
