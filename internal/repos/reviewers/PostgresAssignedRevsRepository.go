package reviewers

import (
	"context"
	"fmt"

	"github.com/salex06/pr-service/internal/database"
)

type PostgresAssignedRevsRepository struct {
	db *database.DB
}

func NewPostgresAssignedRevsRepository(db *database.DB) *PostgresAssignedRevsRepository {
	return &PostgresAssignedRevsRepository{db: db}
}

func (repo *PostgresAssignedRevsRepository) GetAssignedPullRequestIds(ctx context.Context, userId string) ([]string, error) {
	query := `
		SELECT pull_request_id 
		FROM assigned_reviewers
		WHERE user_id = $1;
	`

	rows, err := repo.db.Pool.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get assigned pull requests: %w", err)
	}
	defer rows.Close()

	var prIds []string
	for rows.Next() {
		var prId string
		if err := rows.Scan(prId); err != nil {
			return nil, fmt.Errorf("failed to get team members: %w", err)
		}
		prIds = append(prIds, prId)
	}

	return prIds, nil
}

func (repo *PostgresAssignedRevsRepository) GetAssignedReviewersIds(ctx context.Context, pullRequestId string) ([]string, error) {
	query := `
		SELECT user_id 
		FROM assigned_reviewers
		WHERE pull_request_id = $1;
	`

	rows, err := repo.db.Pool.Query(ctx, query, pullRequestId)
	if err != nil {
		return nil, fmt.Errorf("failed to get assigned pull requests: %w", err)
	}
	defer rows.Close()

	var revIds []string
	for rows.Next() {
		var revId string
		if err := rows.Scan(&revId); err != nil {
			return nil, fmt.Errorf("failed to get team members: %w", err)
		}
		revIds = append(revIds, revId)
	}

	return revIds, nil
}

func (repo *PostgresAssignedRevsRepository) CreateAssignment(ctx context.Context, userId string, prId string) error {
	query := `
		INSERT INTO assigned_reviewers (user_id, pull_request_id)
		VALUES ($1, $2);
	`

	_, err := repo.db.Pool.Exec(ctx, query,
		userId,
		prId,
	)

	if err != nil {
		return fmt.Errorf("failed to create assignment: %w", err)
	}

	return nil
}

func (repo *PostgresAssignedRevsRepository) DeleteAssignment(ctx context.Context, userId string, prId string) error {
	query := `
		DELETE FROM assigned_reviewers
		WHERE user_id = $1 AND pull_request_id = $2
	`

	_, err := repo.db.Pool.Exec(ctx, query,
		userId,
		prId,
	)

	if err != nil {
		return fmt.Errorf("failed to delete assignment: %w", err)
	}

	return nil
}
