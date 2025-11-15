package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/salex06/pr-service/internal/database"
	"github.com/salex06/pr-service/internal/model"
)

type PostgresUserRepository struct {
	db *database.DB
}

func NewPostgresUserRepository(db *database.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (repo *PostgresUserRepository) GetUser(ctx context.Context, userId string) (*model.User, error) {
	query := `
		SELECT user_id, username, team_name, is_active FROM users
		WHERE user_id = $1;
	`

	var user model.User
	err := repo.db.Pool.QueryRow(ctx, query, userId).Scan(
		&user.UserId,
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

func (repo *PostgresUserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users 
		SET username = $1, team_name = $2, is_active = $3
		WHERE user_id = $4;
	`

	result, err := repo.db.Pool.Exec(ctx, query,
		user.Username,
		user.TeamName,
		user.IsActive,
		user.UserId,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (repo *PostgresUserRepository) SaveUser(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (user_id, username, team_name, is_active)
		VALUES ($1, $2, $3, $4);
	`

	_, err := repo.db.Pool.Exec(ctx, query,
		user.UserId,
		user.Username,
		user.TeamName,
		user.IsActive,
	)

	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

func (repo *PostgresUserRepository) UserExists(ctx context.Context, userId string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1);
	`

	err := repo.db.Pool.QueryRow(ctx, query, userId).Scan(&exists)

	return exists, err
}

func (repo *PostgresUserRepository) GetTeamMembers(ctx context.Context, teamName string) ([]*model.User, error) {
	query := `
		SELECT user_id, username, team_name, is_active FROM users
		WHERE team_name=$1; 
	`

	rows, err := repo.db.Pool.Query(ctx, query, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}
	defer rows.Close()

	var members []*model.User
	for rows.Next() {
		var member model.User
		if err := rows.Scan(
			&member.UserId,
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

func (repo *PostgresUserRepository) ChooseReviewers(ctx context.Context, prAuthor *model.User) ([]string, error) {
	query := `
		SELECT user_id FROM users
		WHERE is_active AND team_name=$1 AND user_id != $2
		ORDER BY RANDOM() 
		LIMIT 2;
	`

	rows, err := repo.db.Pool.Query(ctx, query, prAuthor.TeamName, prAuthor.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to choose reviewers: %w", err)
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var revId string
		if err := rows.Scan(&revId); err != nil {
			return nil, fmt.Errorf("failed to choose reviewers: %w", err)
		}
		reviewers = append(reviewers, revId)
	}

	return reviewers, nil
}

func (repo *PostgresUserRepository) ReassignReviewer(ctx context.Context, teamName string, idsExclusionList []string) (*string, error) {
	query := `
		SELECT user_id
		FROM users
		WHERE is_active AND team_name=$1 AND NOT (user_id = ANY($2))
		ORDER BY RANDOM()
		LIMIT 1;
	`

	var userId string
	err := repo.db.Pool.QueryRow(ctx, query, teamName, idsExclusionList).Scan(&userId)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to reassign reviewers: %w", err)
	}

	return &userId, nil
}
