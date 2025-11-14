package team

import "github.com/salex06/pr-service/internal/model"

type TeamRepository interface {
	TeamExists(teamName string) bool
	SaveTeam(team *model.Team) *model.Team
	DeleteMember(teamName string, memberId string)
	GetTeam(teamName string) *model.Team
}
