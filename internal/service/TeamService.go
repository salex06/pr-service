package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/salex06/pr-service/internal/converter"
	"github.com/salex06/pr-service/internal/dto"
	"github.com/salex06/pr-service/internal/model"
	teamRepos "github.com/salex06/pr-service/internal/repos/team"
	userRepos "github.com/salex06/pr-service/internal/repos/user"
)

type TeamService struct {
	teamRepository *teamRepos.TeamRepository
	userRepository *userRepos.UserRepository
}

func NewTeamService(tr *teamRepos.TeamRepository, ur *userRepos.UserRepository) *TeamService {
	return &TeamService{
		teamRepository: tr,
		userRepository: ur,
	}
}

func (ts *TeamService) AddTeam(req *dto.Team) (*dto.Team, *dto.ErrorResponse) {
	teamName := req.TeamName

	if exists, _ := (*ts.teamRepository).TeamExists(context.Background(), teamName); exists {
		return nil, &dto.ErrorResponse{
			Status: http.StatusBadRequest,
			Error: map[string]string{
				"code":    string(dto.TEAM_EXISTS),
				"message": fmt.Sprintf("%s already exists", teamName),
			},
		}
	}

	team := &model.Team{TeamName: teamName}
	(*ts.teamRepository).SaveTeam(context.Background(), team)

	ts.saveMembers(req)

	return &dto.Team{
		TeamName: team.TeamName,
		Members:  req.Members,
	}, nil
}

func (ts *TeamService) saveMembers(req *dto.Team) {
	teamName := req.TeamName

	for _, member := range req.Members {
		if userFromDb, _ := (*ts.userRepository).GetUser(context.Background(), member.UserId); userFromDb != nil {
			// WARN: при создании команды для существующего человека обновляются его поля: teamName, username, isActive

			userFromDb.TeamName = teamName
			userFromDb.Username = member.Username
			userFromDb.IsActive = member.IsActive

			(*ts.userRepository).UpdateUser(context.Background(), userFromDb)
		} else {
			(*ts.userRepository).SaveUser(context.Background(), converter.ConvertTeamMemberToUser(member, teamName))
		}
	}
}

func (ts *TeamService) GetTeam(teamId string) (*dto.Team, *dto.ErrorResponse) {
	if team, _ := (*ts.teamRepository).GetTeam(context.Background(), teamId); team != nil {
		members, _ := (*ts.userRepository).GetTeamMembers(context.Background(), team.TeamName)
		return &dto.Team{
			TeamName: team.TeamName,
			Members:  converter.ConvertUsersToTeamMembers(members),
		}, nil
	}

	return nil, &dto.ErrorResponse{
		Status: http.StatusNotFound,
		Error: map[string]string{
			"code":    string(dto.NOT_FOUND),
			"message": "resource not found",
		},
	}
}
