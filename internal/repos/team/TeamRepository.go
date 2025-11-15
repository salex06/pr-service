package team

import (
	"context"

	"github.com/salex06/pr-service/internal/model"
)

type TeamRepository interface {
	TeamExists(ctx context.Context, teamName string) (bool, error)
	SaveTeam(ctx context.Context, team *model.Team) error
	GetTeam(ctx context.Context, teamName string) (*model.Team, error)
}
