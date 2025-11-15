package team

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/salex06/pr-service/internal/database"
	"github.com/salex06/pr-service/internal/entity"
)

// PostgresTeamRepository представляет собой компонент,
// отвечающий за взаимодействие с БД PostgreSQL, где
// хранится информация о командах
type PostgresTeamRepository struct {
	db *database.DB
}

// NewPostgresTeamRepository конструирует и возвращает объект PostgresTeamRepository
func NewPostgresTeamRepository(db *database.DB) TeamRepository {
	return &PostgresTeamRepository{db: db}
}

// TeamExists выполняет проверку наличия в
// базе данных команды с заданным именем
func (repo *PostgresTeamRepository) TeamExists(ctx context.Context, teamName string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)
	`

	err := repo.db.Pool.QueryRow(ctx, query, teamName).Scan(&exists)

	return exists, err
}

// SaveTeam сохраняет команду в БД
func (repo *PostgresTeamRepository) SaveTeam(ctx context.Context, team *entity.Team) error {
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

// GetTeam выполняет запрос к БД и возвращает
// команду с заданным именем (nil - если не найдена)
func (repo *PostgresTeamRepository) GetTeam(ctx context.Context, teamName string) (*entity.Team, error) {
	query := `
		SELECT team_name FROM teams
		WHERE team_name = $1
	`

	var team entity.Team
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

// GetTeamCount выполняет запрос к БД для
// получения общего количества команд
func (repo *PostgresTeamRepository) GetTeamCount(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM teams
	`

	var count int
	err := repo.db.Pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get team count: %w", err)
	}

	return count, nil
}
