package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/salex06/pr-service/internal/converter"
	"github.com/salex06/pr-service/internal/dto"
	"github.com/salex06/pr-service/internal/entity"
	prRepos "github.com/salex06/pr-service/internal/repos/pr"
	revsRepos "github.com/salex06/pr-service/internal/repos/reviewers"
	teamRepos "github.com/salex06/pr-service/internal/repos/team"
	userRepos "github.com/salex06/pr-service/internal/repos/user"
)

// PullRequestService представляет компонент,
// отвечающий за выполнение бизнес-логики,
// связанной с pull-request`ми
type PullRequestService struct {
	prRepo   *prRepos.PullRequestRepository
	revsRepo *revsRepos.AssignedRevsRepository
	userRepo *userRepos.UserRepository
	teamRepo *teamRepos.TeamRepository
}

// NewPullRequestService конструирует и возвращает объект PullRequestService
func NewPullRequestService(
	prRepo *prRepos.PullRequestRepository,
	revsRepo *revsRepos.AssignedRevsRepository,
	userRepo *userRepos.UserRepository,
	teamRepo *teamRepos.TeamRepository) *PullRequestService {
	return &PullRequestService{
		prRepo:   prRepo,
		revsRepo: revsRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

// CreatePullRequest выполняет открытие нового PR,
// случайным образом назначая до 2-х ревьюеров из
// команды автора PR
func (svc *PullRequestService) CreatePullRequest(req *dto.CreatePullRequest) (*dto.PullRequest, *dto.ErrorResponse) {
	prAuthor, _ := (*svc.userRepo).GetUser(context.Background(), req.AuthorID)
	if prAuthor == nil {
		return nil, &dto.ErrorResponse{
			Status: http.StatusNotFound,
			Error: map[string]string{
				"code":    string(dto.NotFound),
				"message": "resource not found",
			},
		}
	}

	if exists, _ := (*svc.teamRepo).TeamExists(context.Background(), prAuthor.TeamName); !exists {
		return nil, &dto.ErrorResponse{
			Status: http.StatusNotFound,
			Error: map[string]string{
				"code":    string(dto.NotFound),
				"message": "resource not found",
			},
		}
	}

	if exists, _ := (*svc.prRepo).PullRequestExists(context.TODO(), req.PullRequestID); exists {
		return nil, &dto.ErrorResponse{
			Status: http.StatusConflict,
			Error: map[string]string{
				"code":    string(dto.PrExists),
				"message": fmt.Sprintf("%s already exists", req.PullRequestID),
			},
		}
	}

	createTime := time.Now()
	pullRequest := &dto.PullRequest{
		PullRequestID:     req.PullRequestID,
		PullRequestName:   req.PullRequestName,
		AuthorID:          req.AuthorID,
		Status:            entity.OPEN,
		AssignedReviewers: make([]string, 0, dto.MaxAssignedReviewers),
		CreatedAt:         &createTime,
	}
	reviewerIds, err := (*svc.userRepo).ChooseReviewers(context.Background(), prAuthor)
	if err != nil {
		return nil, &dto.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error: map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": fmt.Sprintf("unable choose reviewers: %s", err),
			},
		}
	}

	err = (*svc.prRepo).SavePullRequest(context.Background(), converter.ConvertPrDtoToPrEntity(pullRequest))
	if err != nil {
		return nil, &dto.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error: map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": fmt.Sprintf("unable save PR: %s", err),
			},
		}
	}
	svc.assignReviewers(pullRequest, reviewerIds)

	return pullRequest, nil
}

func (svc *PullRequestService) assignReviewers(pullRequest *dto.PullRequest, reviewers []string) {
	for _, revID := range reviewers {
		err := (*svc.revsRepo).CreateAssignment(context.Background(), revID, pullRequest.PullRequestID)
		if err != nil {
			log.Printf("error occured when creating assignment: %s\n", err)
		}
	}

	pullRequest.AssignedReviewers = append(pullRequest.AssignedReviewers, reviewers...)
}

// MergePullRequest выполняет закрытие PR
// и перевод в статус MERGED
func (svc *PullRequestService) MergePullRequest(req *dto.MergePullRequest) (*dto.PullRequest, *dto.ErrorResponse) {
	pullRequest, _ := (*svc.prRepo).GetPullRequest(context.Background(), req.PullRequestID)
	if pullRequest == nil {
		return nil, &dto.ErrorResponse{
			Status: http.StatusNotFound,
			Error: map[string]string{
				"code":    string(dto.NotFound),
				"message": "resource not found",
			},
		}
	}

	reviewers, _ := (*svc.revsRepo).GetAssignedReviewersIds(context.Background(), pullRequest.PullRequestID)
	if pullRequest.Status == entity.MERGED {
		return converter.ConvertPrToDto(pullRequest, reviewers), nil
	}

	pullRequest.MergedAt = new(time.Time)
	*pullRequest.MergedAt = time.Now()
	pullRequest.Status = entity.MERGED
	err := (*svc.prRepo).UpdatePullRequest(context.Background(), pullRequest)
	if err != nil {
		return nil, &dto.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error: map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": fmt.Sprintf("unable update pull request: %s", err),
			},
		}
	}

	return converter.ConvertPrToDto(pullRequest, reviewers), nil
}

// ReassignPullRequest выполняет переназначение одного сотрудника
// на открытый PR (при наличии активных сотрудников в команде)
func (svc *PullRequestService) ReassignPullRequest(req *dto.ReassignPullRequest) (*dto.ReassignPrResponse, *dto.ErrorResponse) {
	pr, _ := (*svc.prRepo).GetPullRequest(context.Background(), req.PullRequestID)
	userToReplace, _ := (*svc.userRepo).GetUser(context.Background(), req.OldReviewerID)

	if pr == nil || userToReplace == nil {
		return nil, &dto.ErrorResponse{
			Status: http.StatusNotFound,
			Error: map[string]string{
				"code":    string(dto.NotFound),
				"message": "resource not found",
			},
		}
	}

	if pr.Status == entity.MERGED {
		return nil, &dto.ErrorResponse{
			Status: http.StatusConflict,
			Error: map[string]string{
				"code":    string(dto.PrMerged),
				"message": "cannot reassign on merged PR",
			},
		}
	}

	reviewers, _ := (*svc.revsRepo).GetAssignedReviewersIds(context.Background(), pr.PullRequestID)
	if !slices.Contains(reviewers, userToReplace.UserID) {
		return nil, &dto.ErrorResponse{
			Status: http.StatusConflict,
			Error: map[string]string{
				"code":    string(dto.NotAssigned),
				"message": "reviewer is not assigned to this PR",
			},
		}
	}

	idsExclusionList := make([]string, 0, len(reviewers)+1)
	idsExclusionList = append(idsExclusionList, reviewers...)
	idsExclusionList = append(idsExclusionList, pr.AuthorID)

	reassignedReviewerID, _ := (*svc.userRepo).ReassignReviewer(
		context.Background(),
		userToReplace.TeamName,
		idsExclusionList,
	)
	if reassignedReviewerID == nil {
		return nil, &dto.ErrorResponse{
			Status: http.StatusConflict,
			Error: map[string]string{
				"code":    string(dto.NoCandidate),
				"message": "no active replacement candidate in team",
			},
		}
	}

	err := (*svc.revsRepo).DeleteAssignment(context.Background(), userToReplace.UserID, pr.PullRequestID)
	if err != nil {
		return nil, &dto.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error: map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": fmt.Sprintf("unable delete assignment: %s", err),
			},
		}
	}

	err = (*svc.revsRepo).CreateAssignment(context.Background(), *reassignedReviewerID, pr.PullRequestID)
	if err != nil {
		return nil, &dto.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error: map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": fmt.Sprintf("unable create assignment: %s", err),
			},
		}
	}

	reviewers, _ = (*svc.revsRepo).GetAssignedReviewersIds(context.Background(), pr.PullRequestID)
	return converter.ConvertPrToReassigningDto(pr, reviewers, *reassignedReviewerID), nil
}
