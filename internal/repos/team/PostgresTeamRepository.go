package team

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/salex06/pr-service/internal/database"
	"github.com/salex06/pr-service/internal/model"
)

type PostgresTeamRepository struct {
	db *database.DB
}

func NewPostgresTeamRepository(db *database.DB) TeamRepository {
	return &PostgresTeamRepository{db: db}
}

func (repo *PostgresTeamRepository) TeamExists(ctx context.Context, teamName string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)
	`

	err := repo.db.Pool.QueryRow(ctx, query, teamName).Scan(&exists)

	return exists, err
}

func (repo *PostgresTeamRepository) SaveTeam(ctx context.Context, team *model.Team) error {
	query := `
		INSERT INTO teams (team_name)
		VALUES ($1)
	`

	_, err := repo.db.Pool.Exec(ctx, query, team.TeamName)

	if err != nil {
		return fmt.Errorf("failed to save team: %w", err)
	}

	return nil
}

func (repo *PostgresTeamRepository) GetTeam(ctx context.Context, teamName string) (*model.Team, error) {
	query := `
		SELECT team_name FROM teams
		WHERE team_name = $1
	`

	var team model.Team
	err := repo.db.Pool.QueryRow(ctx, query, teamName).Scan(
		&team.TeamName,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	return &team, nil
}
