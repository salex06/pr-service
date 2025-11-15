package team

import (
	"context"

	"github.com/salex06/pr-service/internal/entity"
)

// InMemoryTeamRepository представляет собой компонент,
// который отвечает за взаимодействие с in-memory БД (map),
// где хранится информация о командах
type InMemoryTeamRepository struct {
	storage map[string]*entity.Team
}

// NewInMemoryTeamRepository конструирует и возвращает объект InMemoryTeamRepository
func NewInMemoryTeamRepository() *InMemoryTeamRepository {
	return &InMemoryTeamRepository{
		storage: make(map[string]*entity.Team),
	}
}

// TeamExists проверяет, содержится ли команда с заданным именем в map
func (db *InMemoryTeamRepository) TeamExists(ctx context.Context, teamName string) (bool, error) {
	_, ok := db.storage[teamName]
	return ok, nil
}

// SaveTeam сохраняет команду в map по заданному имени
func (db *InMemoryTeamRepository) SaveTeam(ctx context.Context, team *entity.Team) error {
	db.storage[team.TeamName] = team

	return nil
}

// GetTeam возвращает команду с заданным именем (nil - если не найдена)
func (db *InMemoryTeamRepository) GetTeam(ctx context.Context, teamName string) (*entity.Team, error) {
	return db.storage[teamName], nil
}
