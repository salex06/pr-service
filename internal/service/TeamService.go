package service

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/salex06/pr-service/internal/converter"
	"github.com/salex06/pr-service/internal/dto"
	"github.com/salex06/pr-service/internal/entity"
	teamRepos "github.com/salex06/pr-service/internal/repos/team"
	userRepos "github.com/salex06/pr-service/internal/repos/user"
)

// TeamService представляет компонент, отвечающий за
// выполнение бизнес-логики, связанной с командами -
// группами пользователей с уникальным именем
type TeamService struct {
	teamRepository *teamRepos.TeamRepository
	userRepository *userRepos.UserRepository
}

// NewTeamService конструирует и возвращает объект TeamService
func NewTeamService(tr *teamRepos.TeamRepository, ur *userRepos.UserRepository) *TeamService {
	return &TeamService{
		teamRepository: tr,
		userRepository: ur,
	}
}

// AddTeam выполняет сохранение команды и её представителей
func (ts *TeamService) AddTeam(req *dto.Team) (*dto.Team, *dto.ErrorResponse) {
	teamName := req.TeamName

	if exists, _ := (*ts.teamRepository).TeamExists(context.Background(), teamName); exists {
		return nil, &dto.ErrorResponse{
			Status: http.StatusBadRequest,
			Error: map[string]string{
				"code":    string(dto.TeamExists),
				"message": fmt.Sprintf("%s already exists", teamName),
			},
		}
	}

	team := &entity.Team{TeamName: teamName}
	err := (*ts.teamRepository).SaveTeam(context.Background(), team)
	if err != nil {
		return nil, &dto.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error: map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": fmt.Sprintf("unable save team: %s", err),
			},
		}
	}

	ts.saveMembers(req)

	return &dto.Team{
		TeamName: team.TeamName,
		Members:  req.Members,
	}, nil
}

func (ts *TeamService) saveMembers(req *dto.Team) {
	teamName := req.TeamName

	for _, member := range req.Members {
		if userFromDB, _ := (*ts.userRepository).GetUser(context.Background(), member.UserID); userFromDB != nil {
			// WARN: при создании команды для существующего человека обновляются его поля: teamName, username, isActive

			userFromDB.TeamName = teamName
			userFromDB.Username = member.Username
			userFromDB.IsActive = member.IsActive

			err := (*ts.userRepository).UpdateUser(context.Background(), userFromDB)
			if err != nil {
				log.Printf("error occured when updating user: %s\n", err)
			}
		} else {
			err := (*ts.userRepository).SaveUser(context.Background(), converter.ConvertTeamMemberToUser(member, teamName))
			if err != nil {
				log.Printf("error occured when saving user: %s\n", err)
			}
		}
	}
}

// GetTeam возвращает объект команды,
// имеющей идентификатор teamID
func (ts *TeamService) GetTeam(teamID string) (*dto.Team, *dto.ErrorResponse) {
	if team, _ := (*ts.teamRepository).GetTeam(context.Background(), teamID); team != nil {
		members, _ := (*ts.userRepository).GetTeamMembers(context.Background(), team.TeamName)
		return &dto.Team{
			TeamName: team.TeamName,
			Members:  converter.ConvertUsersToTeamMembers(members),
		}, nil
	}

	return nil, &dto.ErrorResponse{
		Status: http.StatusNotFound,
		Error: map[string]string{
			"code":    string(dto.NotFound),
			"message": "resource not found",
		},
	}
}
