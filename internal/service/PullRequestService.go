package service

import (
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
	prAuthor := (*svc.userRepo).GetUser(req.AuthorId)
	if prAuthor == nil || !(*svc.teamRepo).TeamExists(prAuthor.TeamName) {
		return nil, &dto.ErrorResponse{
			Status: http.StatusNotFound,
			Error: map[string]string{
				"code":    string(dto.NOT_FOUND),
				"message": "resource not found",
			},
		}
	}

	if (*svc.prRepo).PullRequestExists(req.PullRequestId) {
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
	reviewerIds := (*svc.userRepo).ChooseReviewers(prAuthor)
	svc.assignReviewers(pullRequest, reviewerIds)
	(*svc.prRepo).SavePullRequest(svc.convertDtoToPr(pullRequest))

	return pullRequest, nil
}

func (svc *PullRequestService) assignReviewers(pullRequest *dto.PullRequest, reviewers []string) {
	for _, revId := range reviewers {
		(*svc.revsRepo).CreateAssignment(revId, pullRequest.PullRequestId)
	}

	pullRequest.AssignedReviewers = append(pullRequest.AssignedReviewers, reviewers...)
}

func (svc *PullRequestService) convertDtoToPr(pr *dto.PullRequest) *model.PullRequest {
	return &model.PullRequest{
		PullRequestId:   pr.PullRequestId,
		PullRequestName: pr.PullRequestName,
		AuthorId:        pr.AuthorId,
		Status:          pr.Status,
		CreatedAt:       *pr.CreatedAt,
		MergedAt:        *pr.MergedAt,
	}
}

func (svc *PullRequestService) MergePullRequest(req *dto.MergePullRequest) (*dto.PullRequest, *dto.ErrorResponse) {
	pullRequest := (*svc.prRepo).GetPullRequest(req.PullRequestId)
	if pullRequest == nil {
		return nil, &dto.ErrorResponse{
			Status: http.StatusNotFound,
			Error: map[string]string{
				"code":    string(dto.NOT_FOUND),
				"message": "resource not found",
			},
		}
	}

	reviewers := (*svc.revsRepo).GetAssignedReviewersIds(pullRequest.PullRequestId)
	if pullRequest.Status == model.MERGED {
		return svc.convertPrToDto(pullRequest, reviewers), nil
	}

	pullRequest.MergedAt = time.Now()
	(*svc.prRepo).UpdatePullRequest(pullRequest)

	return svc.convertPrToDto(pullRequest, reviewers), nil
}

func (svc *PullRequestService) convertPrToDto(pr *model.PullRequest, reviewers []string) *dto.PullRequest {
	return &dto.PullRequest{
		PullRequestId:     pr.PullRequestId,
		PullRequestName:   pr.PullRequestName,
		AuthorId:          pr.AuthorId,
		Status:            pr.Status,
		AssignedReviewers: reviewers,
		CreatedAt:         &pr.CreatedAt,
		MergedAt:          &pr.MergedAt,
	}
}

func (svc *PullRequestService) ReassignPullRequest(req *dto.ReassignPullRequest) (*dto.ReassignPrResponse, *dto.ErrorResponse) {
	pr := (*svc.prRepo).GetPullRequest(req.PullRequestId)
	userToReplace := (*svc.userRepo).GetUser(req.OldReviewerId)

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

	reviewers := (*svc.revsRepo).GetAssignedReviewersIds(pr.PullRequestId)
	if !slices.Contains(reviewers, userToReplace.UserId) {
		return nil, &dto.ErrorResponse{
			Status: http.StatusConflict,
			Error: map[string]string{
				"code":    string(dto.NOT_ASSIGNED),
				"message": "reviewer is not assigned to this PR",
			},
		}
	}

	reassignedReviewerId := (*svc.userRepo).ReassignReviewer(userToReplace.TeamName, []string{pr.AuthorId, userToReplace.UserId})
	if reassignedReviewerId == nil {
		return nil, &dto.ErrorResponse{
			Status: http.StatusConflict,
			Error: map[string]string{
				"code":    string(dto.NO_CANDIDATE),
				"message": "no active replacement candidate in team",
			},
		}
	}

	(*svc.revsRepo).DeleteAssignment(userToReplace.UserId, pr.PullRequestId)
	(*svc.revsRepo).CreateAssignment(*reassignedReviewerId, pr.PullRequestId)

	reviewers = (*svc.revsRepo).GetAssignedReviewersIds(pr.PullRequestId)
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
			CreatedAt:         &pr.CreatedAt,
			MergedAt:          &pr.MergedAt,
		},
		ReplacedBy: replacedBy,
	}
}
