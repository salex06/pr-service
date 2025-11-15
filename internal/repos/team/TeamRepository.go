// Package team - пакет с репозиториями, отвечающими за взаимодействие с БД,
// где хранится информация о командах
package team

import (
	"context"

	"github.com/salex06/pr-service/internal/entity"
)

// TeamRepository представляет интерфейс взаимодействия
// с базой данных, где хранится информация о командах
type TeamRepository interface {
	TeamExists(ctx context.Context, teamName string) (bool, error)
	SaveTeam(ctx context.Context, team *entity.Team) error
	GetTeam(ctx context.Context, teamName string) (*entity.Team, error)

	GetTeamCount(ctx context.Context) (int, error)
}
