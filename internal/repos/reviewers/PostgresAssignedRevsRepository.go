package reviewers

import (
	"context"
	"fmt"

	"github.com/salex06/pr-service/internal/database"
)

// PostgresAssignedRevsRepository представляет собой компонент,
// отвечающий за взаимодействие с БД PostgreSQL, где хранится
// информация о назначениях сотрудников на PR's
type PostgresAssignedRevsRepository struct {
	db *database.DB
}

// NewPostgresAssignedRevsRepository конструирует и возвращает объект PostgresAssignedRevsRepository
func NewPostgresAssignedRevsRepository(db *database.DB) AssignedRevsRepository {
	return &PostgresAssignedRevsRepository{db: db}
}

// GetAssignedPullRequestIds выполняет запрос к БД и возвращает
// слайс идентификаторов PR`s, на которые назначен сотрудник
func (repo *PostgresAssignedRevsRepository) GetAssignedPullRequestIds(ctx context.Context, userID string) ([]string, error) {
	query := `
		SELECT pull_request_id 
		FROM assigned_reviewers
		WHERE user_id = $1;
	`

	rows, err := repo.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assigned pull requests: %w", err)
	}
	defer rows.Close()

	var prIds []string
	for rows.Next() {
		var prID string
		if err := rows.Scan(&prID); err != nil {
			return nil, fmt.Errorf("failed to get team members: %w", err)
		}
		prIds = append(prIds, prID)
	}

	return prIds, nil
}

// GetAssignedReviewersIds выполняет запрос к БД и возвращает
// слайс идентификаторов сотрудников, назначенных на данный PR
func (repo *PostgresAssignedRevsRepository) GetAssignedReviewersIds(ctx context.Context, pullRequestID string) ([]string, error) {
	query := `
		SELECT user_id 
		FROM assigned_reviewers
		WHERE pull_request_id = $1;
	`

	rows, err := repo.db.Pool.Query(ctx, query, pullRequestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assigned pull requests: %w", err)
	}
	defer rows.Close()

	var revIds []string
	for rows.Next() {
		var revID string
		if err := rows.Scan(&revID); err != nil {
			return nil, fmt.Errorf("failed to get team members: %w", err)
		}
		revIds = append(revIds, revID)
	}

	return revIds, nil
}

// CreateAssignment выполняет запрос к БД для
// назначения сотрудника с идентификатором userID
// на PR с идентификатором prID
func (repo *PostgresAssignedRevsRepository) CreateAssignment(ctx context.Context, userID, prID string) error {
	query := `
		INSERT INTO assigned_reviewers (user_id, pull_request_id)
		VALUES ($1, $2);
	`

	_, err := repo.db.Pool.Exec(ctx, query,
		userID,
		prID,
	)

	if err != nil {
		return fmt.Errorf("failed to create assignment: %w", err)
	}

	return nil
}

// DeleteAssignment выполняет запрос к БД для удаления
// назначения сотрудника с идентификатором userID на PR
// с идентификатором prID
func (repo *PostgresAssignedRevsRepository) DeleteAssignment(ctx context.Context, userID, prID string) error {
	query := `
		DELETE FROM assigned_reviewers
		WHERE user_id = $1 AND pull_request_id = $2
	`

	_, err := repo.db.Pool.Exec(ctx, query,
		userID,
		prID,
	)

	if err != nil {
		return fmt.Errorf("failed to delete assignment: %w", err)
	}

	return nil
}
