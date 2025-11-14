package service

import (
	"fmt"
	"net/http"

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

	if (*ts.teamRepository).TeamExists(teamName) {
		return nil, &dto.ErrorResponse{
			Status: http.StatusBadRequest,
			Error: map[string]string{
				"code":    string(dto.TEAM_EXISTS),
				"message": fmt.Sprintf("%s already exists", teamName),
			},
		}
	}

	team := &model.Team{
		TeamName: teamName,
		Members:  make([]*model.User, 0),
	}
	for _, member := range req.Members {
		var savedUser *model.User
		if userFromDb := (*ts.userRepository).GetUser(member.UserId); userFromDb != nil {
			(*ts.teamRepository).DeleteMember(userFromDb.TeamName, userFromDb.UserId)
			userFromDb.TeamName = teamName
			//TODO: также необходимо обновлять флаг состояния is_active
			savedUser = (*ts.userRepository).UpdateUser(userFromDb)
		} else {
			savedUser = (*ts.userRepository).SaveUser(ts.convertTeamMemberToUser(member, teamName))
		}

		team.Members = append(team.Members, savedUser)
	}

	(*ts.teamRepository).SaveTeam(team)
	return &dto.Team{
		TeamName: team.TeamName,
		Members:  ts.convertUsersToTeamMembers(team.Members),
	}, nil
}

func (ts *TeamService) convertTeamMemberToUser(member *dto.TeamMember, teamName string) *model.User {
	return &model.User{
		UserId:   member.UserId,
		Username: member.Username,
		TeamName: teamName,
		IsActive: member.IsActive,
	}
}

func (ts *TeamService) convertUsersToTeamMembers(users []*model.User) []*dto.TeamMember {
	converted := make([]*dto.TeamMember, 0, len(users))
	for _, user := range users {
		converted = append(converted, ts.convertUserToTeamMember(user))
	}

	return converted
}

func (ts *TeamService) convertUserToTeamMember(user *model.User) *dto.TeamMember {
	return &dto.TeamMember{
		UserId:   user.UserId,
		Username: user.Username,
		IsActive: user.IsActive,
	}
}

func (ts *TeamService) GetTeam(teamId string) (*dto.Team, *dto.ErrorResponse) {
	if team := (*ts.teamRepository).GetTeam(teamId); team != nil {
		return &dto.Team{
			TeamName: team.TeamName,
			Members:  ts.convertUsersToTeamMembers(team.Members),
		}, nil
	}

	return nil, &dto.ErrorResponse{
		Status: 404,
		Error: map[string]string{
			"code":    string(dto.NOT_FOUND),
			"message": "resource not found",
		},
	}
}
