package team

import (
	"context"

	"github.com/salex06/pr-service/internal/model"
)

type InMemoryTeamRepository struct {
	storage map[string]*model.Team
}

func NewInMemoryTeamRepository() *InMemoryTeamRepository {
	return &InMemoryTeamRepository{
		storage: make(map[string]*model.Team),
	}
}

func (db *InMemoryTeamRepository) TeamExists(ctx context.Context, teamName string) (bool, error) {
	_, ok := db.storage[teamName]
	return ok, nil
}

func (db *InMemoryTeamRepository) SaveTeam(ctx context.Context, team *model.Team) error {
	db.storage[team.TeamName] = team

	return nil
}

func (db *InMemoryTeamRepository) GetTeam(ctx context.Context, teamName string) (*model.Team, error) {
	return db.storage[teamName], nil
}
