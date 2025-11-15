package service

import (
	"context"
	"net/http"

	"github.com/salex06/pr-service/internal/dto"
	prRepos "github.com/salex06/pr-service/internal/repos/pr"
	revsRepos "github.com/salex06/pr-service/internal/repos/reviewers"
	teamRepos "github.com/salex06/pr-service/internal/repos/team"
	userRepos "github.com/salex06/pr-service/internal/repos/user"
)

// StatsService представляет собой компонент,
// который отвечает за выполнение бизнес-логики,
// связанной со сбором статистики работы приложения
type StatsService struct {
	prRepo   *prRepos.PullRequestRepository
	revsRepo *revsRepos.AssignedRevsRepository
	userRepo *userRepos.UserRepository
	teamRepo *teamRepos.TeamRepository
}

// NewStatsService контструирует и возвращает объект StatsService
func NewStatsService(
	prRepo *prRepos.PullRequestRepository,
	revsRepo *revsRepos.AssignedRevsRepository,
	userRepo *userRepos.UserRepository,
	teamRepo *teamRepos.TeamRepository) *StatsService {
	return &StatsService{
		prRepo:   prRepo,
		revsRepo: revsRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

// GetStat возвращает базовую статистику по состоянию приложения
// (в данный момент только бизнес-метрики)
func (svc *StatsService) GetStat() (*dto.AppStat, *dto.ErrorResponse) {
	var stat dto.AppStat

	userCountInfo, err := svc.getUserCountInfo()
	if err != nil {
		return nil, &dto.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error: map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": "error occured when getting user count info",
			},
		}
	}

	teamCount, err := svc.getTeamCount()
	if err != nil {
		return nil, &dto.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error: map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": "error occured when getting team count info",
			},
		}
	}

	prCountInfo, err := svc.getPrCountInfo()
	if err != nil {
		return nil, &dto.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error: map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": "error occured when getting PR count info",
			},
		}
	}

	userCountByTeams, err := svc.getUserCountGroupedByTeams()
	if err != nil {
		return nil, &dto.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error: map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": "error occured when getting user count by teams info",
			},
		}
	}

	assignmentsCountByUser, err := svc.getAssignmentsCountByUsers()
	if err != nil {
		return nil, &dto.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error: map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": "error occured when getting assignments count by user info",
			},
		}
	}

	stat.TotalUsersCount = userCountInfo["total"]
	stat.ActiveUsersCount = userCountInfo["active"]
	stat.TotalTeamsCount = teamCount
	stat.OpenedPRCount = prCountInfo["open"]
	stat.MergedPRCount = prCountInfo["merged"]
	stat.UserCountByTeam = userCountByTeams
	stat.AssignmentsCountByUser = assignmentsCountByUser

	return &stat, nil
}

func (svc *StatsService) getUserCountInfo() (map[string]int, error) {
	totalUserCount, err := (*svc.userRepo).GetTotalUserCount(context.Background())
	if err != nil {
		return nil, err
	}

	activeUserCount, err := (*svc.userRepo).GetActiveUserCount(context.Background())
	if err != nil {
		return nil, err
	}

	return map[string]int{
		"total":  totalUserCount,
		"active": activeUserCount,
	}, nil
}

func (svc *StatsService) getTeamCount() (int, error) {
	teamCount, err := (*svc.teamRepo).GetTeamCount(context.Background())
	if err != nil {
		return 0, err
	}

	return teamCount, nil
}

func (svc *StatsService) getPrCountInfo() (map[string]int, error) {
	openPrCount, err := (*svc.prRepo).GetOpenedPullRequestCount(context.Background())
	if err != nil {
		return nil, err
	}

	mergedPrCount, err := (*svc.prRepo).GetMergedPullRequestCount(context.Background())
	if err != nil {
		return nil, err
	}

	return map[string]int{
		"open":   openPrCount,
		"merged": mergedPrCount,
	}, nil
}

func (svc *StatsService) getUserCountGroupedByTeams() ([]*dto.TeamSize, error) {
	count, err := (*svc.userRepo).GetUserCountByTeam(context.Background())
	if err != nil {
		return nil, err
	}

	return count, nil
}

func (svc *StatsService) getAssignmentsCountByUsers() ([]*dto.AssignmentsByUser, error) {
	count, err := (*svc.revsRepo).GetAssignmentsCountByReviewerID(context.Background())
	if err != nil {
		return nil, err
	}

	return count, nil
}
