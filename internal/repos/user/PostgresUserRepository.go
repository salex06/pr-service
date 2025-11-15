package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/salex06/pr-service/internal/database"
	"github.com/salex06/pr-service/internal/dto"
	"github.com/salex06/pr-service/internal/entity"
)

// PostgresUserRepository представляет собой компонент,
// отвечающей за взаимодействие с БД, где хранится информация
// о пользователях
type PostgresUserRepository struct {
	db *database.DB
}

// NewPostgresUserRepository конструирует и возвращает объект PostgresUserRepository
func NewPostgresUserRepository(db *database.DB) UserRepository {
	return &PostgresUserRepository{db: db}
}

// GetUser возвращает пользователя с заданным userID (если не найден - nil)
func (repo *PostgresUserRepository) GetUser(ctx context.Context, userID string) (*entity.User, error) {
	query := `
		SELECT user_id, username, team_name, is_active FROM users
		WHERE user_id = $1;
	`

	var user entity.User
	err := repo.db.Pool.QueryRow(ctx, query, userID).Scan(
		&user.UserID,
		&user.Username,
		&user.TeamName,
		&user.IsActive,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// UpdateUser выполняет запрос к БД для обновления
// изменяемой информации о пользователе
func (repo *PostgresUserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users 
		SET username = $1, team_name = $2, is_active = $3
		WHERE user_id = $4;
	`

	result, err := repo.db.Pool.Exec(ctx, query,
		user.Username,
		user.TeamName,
		user.IsActive,
		user.UserID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}

// SaveUser сохраняет пользователя в БД
func (repo *PostgresUserRepository) SaveUser(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (user_id, username, team_name, is_active)
		VALUES ($1, $2, $3, $4);
	`

	_, err := repo.db.Pool.Exec(ctx, query,
		user.UserID,
		user.Username,
		user.TeamName,
		user.IsActive,
	)

	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

// UserExists проверяет, существует ли в БД
// пользователь с заданным id, и возвращает результат
func (repo *PostgresUserRepository) UserExists(ctx context.Context, userID string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1);
	`

	err := repo.db.Pool.QueryRow(ctx, query, userID).Scan(&exists)

	return exists, err
}

// GetTotalUserCount обращается к БД и возвращает
// общее количество пользователей
func (repo *PostgresUserRepository) GetTotalUserCount(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM users
	`

	var count int
	err := repo.db.Pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get users count: %w", err)
	}

	return count, nil
}

// GetActiveUserCount обращается к БД и возвращает
// количество активных пользователей
func (repo *PostgresUserRepository) GetActiveUserCount(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM users
		WHERE is_active
	`

	var count int
	err := repo.db.Pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get active users count: %w", err)
	}

	return count, nil
}

// GetTeamMembers выполняет запрос к БД и возвращает
// сотрудников, которые являются членами заданной команды
func (repo *PostgresUserRepository) GetTeamMembers(ctx context.Context, teamName string) ([]*entity.User, error) {
	query := `
		SELECT user_id, username, team_name, is_active FROM users
		WHERE team_name=$1; 
	`

	rows, err := repo.db.Pool.Query(ctx, query, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}
	defer rows.Close()

	var members []*entity.User
	for rows.Next() {
		var member entity.User
		if err := rows.Scan(
			&member.UserID,
			&member.Username,
			&member.TeamName,
			&member.IsActive,
		); err != nil {
			return nil, fmt.Errorf("failed to get team members: %w", err)
		}
		members = append(members, &member)
	}

	return members, nil
}

// GetUserCountByTeam выполняет запрос к БД и возвращает
// количество пользователей в каждой команде
func (repo *PostgresUserRepository) GetUserCountByTeam(ctx context.Context) ([]*dto.TeamSize, error) {
	query := `
		SELECT team_name, COUNT(*)
		FROM users
		GROUP BY team_name;
	`
	rows, err := repo.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to group users by teams: %w", err)
	}
	defer rows.Close()

	teamSizes := make([]*dto.TeamSize, 0)
	for rows.Next() {
		var teamSize dto.TeamSize
		if err := rows.Scan(&teamSize.TeamName, &teamSize.UserCount); err != nil {
			return nil, fmt.Errorf("failed to group users by teams: %w", err)
		}
		teamSizes = append(teamSizes, &teamSize)
	}

	return teamSizes, nil
}

// ChooseReviewers выполняет запрос к БД и возвращает
// до 2-х сотрудников, которые могут быть назначены на PR
func (repo *PostgresUserRepository) ChooseReviewers(ctx context.Context, prAuthor *entity.User) ([]string, error) {
	query := `
		SELECT user_id FROM users
		WHERE is_active AND team_name=$1 AND user_id != $2
		ORDER BY RANDOM() 
		LIMIT $3;
	`

	rows, err := repo.db.Pool.Query(ctx, query, prAuthor.TeamName, prAuthor.UserID, dto.MaxAssignedReviewers)
	if err != nil {
		return nil, fmt.Errorf("failed to choose reviewers: %w", err)
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var revID string
		if err := rows.Scan(&revID); err != nil {
			return nil, fmt.Errorf("failed to choose reviewers: %w", err)
		}
		reviewers = append(reviewers, revID)
	}

	return reviewers, nil
}

// ReassignReviewer выполняет запрос к БД и возвращает
// случайного активного сотрудника из команды автора PR,
// который может быть назначен на PR вместо прежнего сотрудника
func (repo *PostgresUserRepository) ReassignReviewer(ctx context.Context, teamName string, idsExclusionList []string) (*string, error) {
	query := `
		SELECT user_id
		FROM users
		WHERE is_active AND team_name=$1 AND NOT (user_id = ANY($2))
		ORDER BY RANDOM()
		LIMIT 1;
	`

	var userID string
	err := repo.db.Pool.QueryRow(ctx, query, teamName, idsExclusionList).Scan(&userID)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to reassign reviewers: %w", err)
	}

	return &userID, nil
}
