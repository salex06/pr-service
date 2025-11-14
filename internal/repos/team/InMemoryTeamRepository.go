package team

import "github.com/salex06/pr-service/internal/model"

type InMemoryTeamRepository struct {
	storage map[string]*model.Team
}

func NewInMemoryTeamRepository() *InMemoryTeamRepository {
	return &InMemoryTeamRepository{
		storage: make(map[string]*model.Team),
	}
}

func (db *InMemoryTeamRepository) TeamExists(teamName string) bool {
	_, ok := db.storage[teamName]
	return ok
}

func (db *InMemoryTeamRepository) SaveTeam(team *model.Team) *model.Team {
	db.storage[team.TeamName] = team

	return db.storage[team.TeamName]
}

func (db *InMemoryTeamRepository) DeleteMember(teamName string, userId string) {
	for k, v := range db.storage {
		if k == teamName {
			sl := make([]*model.User, 0, len(v.Members))
			for j := range v.Members {
				if v.Members[j].UserId != userId {
					sl = append(sl, v.Members[j])
				}
			}
			v.Members = sl
		}
	}
}

func (db *InMemoryTeamRepository) GetTeam(teamName string) *model.Team {
	return db.storage[teamName]
}
