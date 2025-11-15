package service

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/salex06/pr-service/internal/dto"
	"github.com/salex06/pr-service/internal/model"
	prRepos "github.com/salex06/pr-service/internal/repos/pr"
	revsRepos "github.com/salex06/pr-service/internal/repos/reviewers"
	teamRepos "github.com/salex06/pr-service/internal/repos/team"
	userRepos "github.com/salex06/pr-service/internal/repos/user"
)

type PullRequestService struct {
	prRepo   *prRepos.PullRequestRepository
	revsRepo *revsRepos.AssignedRevsRepository
	userRepo *userRepos.UserRepository
	teamRepo *teamRepos.TeamRepository
}

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

func (svc *PullRequestService) CreatePullRequest(req *dto.CreatePullRequest) (*dto.PullRequest, *dto.ErrorResponse) {
	prAuthor, _ := (*svc.userRepo).GetUser(context.Background(), req.AuthorId)
	if prAuthor == nil {
		return nil, &dto.ErrorResponse{
			Status: http.StatusNotFound,
			Error: map[string]string{
				"code":    string(dto.NOT_FOUND),
				"message": "resource not found",
			},
		}
	}
	if exists, _ := (*svc.teamRepo).TeamExists(context.Background(), prAuthor.TeamName); !exists {
		return nil, &dto.ErrorResponse{
			Status: http.StatusNotFound,
			Error: map[string]string{
				"code":    string(dto.NOT_FOUND),
				"message": "resource not found",
			},
		}
	}

	if exists, _ := (*svc.prRepo).PullRequestExists(context.TODO(), req.PullRequestId); exists {
		return nil, &dto.ErrorResponse{
			Status: http.StatusConflict,
			Error: map[string]string{
				"code":    string(dto.PR_EXISTS),
				"message": fmt.Sprintf("%s already exists", req.PullRequestId),
			},
		}
	}

	pullRequest := &dto.PullRequest{
		PullRequestId:     req.PullRequestId,
		PullRequestName:   req.PullRequestName,
		AuthorId:          req.AuthorId,
		Status:            model.OPEN,
		AssignedReviewers: make([]string, 0, dto.MAX_ASSIGNED_REVIEWERS),
	}
	reviewerIds, _ := (*svc.userRepo).ChooseReviewers(context.Background(), prAuthor)
	(*svc.prRepo).SavePullRequest(context.Background(), svc.convertDtoToPr(pullRequest))
	svc.assignReviewers(pullRequest, reviewerIds)

	return pullRequest, nil
}

func (svc *PullRequestService) assignReviewers(pullRequest *dto.PullRequest, reviewers []string) {
	for _, revId := range reviewers {
		(*svc.revsRepo).CreateAssignment(context.Background(), revId, pullRequest.PullRequestId)
	}

	pullRequest.AssignedReviewers = append(pullRequest.AssignedReviewers, reviewers...)
}

func (svc *PullRequestService) convertDtoToPr(pr *dto.PullRequest) *model.PullRequest {
	return &model.PullRequest{
		PullRequestId:   pr.PullRequestId,
		PullRequestName: pr.PullRequestName,
		AuthorId:        pr.AuthorId,
		Status:          pr.Status,
		CreatedAt:       pr.CreatedAt,
		MergedAt:        pr.MergedAt,
	}
}

func (svc *PullRequestService) MergePullRequest(req *dto.MergePullRequest) (*dto.PullRequest, *dto.ErrorResponse) {
	pullRequest, _ := (*svc.prRepo).GetPullRequest(context.Background(), req.PullRequestId)
	if pullRequest == nil {
		return nil, &dto.ErrorResponse{
			Status: http.StatusNotFound,
			Error: map[string]string{
				"code":    string(dto.NOT_FOUND),
				"message": "resource not found",
			},
		}
	}

	reviewers, _ := (*svc.revsRepo).GetAssignedReviewersIds(context.Background(), pullRequest.PullRequestId)
	if pullRequest.Status == model.MERGED {
		return svc.convertPrToDto(pullRequest, reviewers), nil
	}

	pullRequest.MergedAt = new(time.Time)
	*pullRequest.MergedAt = time.Now()
	pullRequest.Status = model.MERGED
	(*svc.prRepo).UpdatePullRequest(context.Background(), pullRequest)

	return svc.convertPrToDto(pullRequest, reviewers), nil
}

func (svc *PullRequestService) convertPrToDto(pr *model.PullRequest, reviewers []string) *dto.PullRequest {
	return &dto.PullRequest{
		PullRequestId:     pr.PullRequestId,
		PullRequestName:   pr.PullRequestName,
		AuthorId:          pr.AuthorId,
		Status:            pr.Status,
		AssignedReviewers: reviewers,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}

func (svc *PullRequestService) ReassignPullRequest(req *dto.ReassignPullRequest) (*dto.ReassignPrResponse, *dto.ErrorResponse) {
	pr, _ := (*svc.prRepo).GetPullRequest(context.Background(), req.PullRequestId)
	userToReplace, _ := (*svc.userRepo).GetUser(context.Background(), req.OldReviewerId)

	if pr == nil || userToReplace == nil {
		return nil, &dto.ErrorResponse{
			Status: http.StatusNotFound,
			Error: map[string]string{
				"code":    string(dto.NOT_FOUND),
				"message": "resource not found",
			},
		}
	}

	if pr.Status == model.MERGED {
		return nil, &dto.ErrorResponse{
			Status: http.StatusConflict,
			Error: map[string]string{
				"code":    string(dto.PR_MERGED),
				"message": "cannot reassign on merged PR",
			},
		}
	}

	reviewers, _ := (*svc.revsRepo).GetAssignedReviewersIds(context.Background(), pr.PullRequestId)
	if !slices.Contains(reviewers, userToReplace.UserId) {
		return nil, &dto.ErrorResponse{
			Status: http.StatusConflict,
			Error: map[string]string{
				"code":    string(dto.NOT_ASSIGNED),
				"message": "reviewer is not assigned to this PR",
			},
		}
	}

	idsExclusionList := make([]string, 0, len(reviewers)+1)
	idsExclusionList = append(idsExclusionList, reviewers...)
	idsExclusionList = append(idsExclusionList, pr.AuthorId)

	reassignedReviewerId, _ := (*svc.userRepo).ReassignReviewer(
		context.Background(),
		userToReplace.TeamName,
		idsExclusionList,
	)
	if reassignedReviewerId == nil {
		return nil, &dto.ErrorResponse{
			Status: http.StatusConflict,
			Error: map[string]string{
				"code":    string(dto.NO_CANDIDATE),
				"message": "no active replacement candidate in team",
			},
		}
	}

	(*svc.revsRepo).DeleteAssignment(context.Background(), userToReplace.UserId, pr.PullRequestId)
	(*svc.revsRepo).CreateAssignment(context.Background(), *reassignedReviewerId, pr.PullRequestId)

	reviewers, _ = (*svc.revsRepo).GetAssignedReviewersIds(context.Background(), pr.PullRequestId)
	return svc.convertPrToReassingDto(pr, reviewers, *reassignedReviewerId), nil
}

func (svc *PullRequestService) convertPrToReassingDto(pr *model.PullRequest, reviewers []string, replacedBy string) *dto.ReassignPrResponse {
	return &dto.ReassignPrResponse{
		Pr: dto.PullRequest{
			PullRequestId:     pr.PullRequestId,
			PullRequestName:   pr.PullRequestName,
			AuthorId:          pr.AuthorId,
			Status:            pr.Status,
			AssignedReviewers: reviewers,
			CreatedAt:         pr.CreatedAt,
			MergedAt:          pr.MergedAt,
		},
		ReplacedBy: replacedBy,
	}
}
