package pr

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/salex06/pr-service/internal/database"
	"github.com/salex06/pr-service/internal/entity"
)

// PostgresPullRequestRepository представляет собой компонент,
// отвечающий за взаимодействие с БД PostgreSQL, где
// содержится информация о PR's
type PostgresPullRequestRepository struct {
	db *database.DB
}

// NewPostgresPullRequestRepository конструирует и возвращает объект PostgresPullRequestRepository
func NewPostgresPullRequestRepository(db *database.DB) PullRequestRepository {
	return &PostgresPullRequestRepository{db: db}
}

// PullRequestExists выполняет запрос для проверки
// наличия в БД PR с заданным идентификатором
func (repo *PostgresPullRequestRepository) PullRequestExists(ctx context.Context, prIВ string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1)
	`

	err := repo.db.Pool.QueryRow(ctx, query, prIВ).Scan(&exists)

	return exists, err
}

// GetPullRequest выполняет запрос к БД для получения
// PR с заданным идентификатором (nil - если не найден)
func (repo *PostgresPullRequestRepository) GetPullRequest(ctx context.Context, prID string) (*entity.PullRequest, error) {
	query := `
		SELECT pull_request_id, pull_request_name, author_id, pr_status, created_at, merged_at 
		FROM pull_requests
		WHERE pull_request_id = $1
	`

	var pr entity.PullRequest
	err := repo.db.Pool.QueryRow(ctx, query, prID).Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get pull request: %w", err)
	}

	return &pr, nil
}

// GetPullRequests возвращает набор объектов PR's по заданному набору идентификаторов
func (repo *PostgresPullRequestRepository) GetPullRequests(ctx context.Context, prIds []string) ([]*entity.PullRequest, error) {
	prs := make([]*entity.PullRequest, 0, len(prIds))
	for _, id := range prIds {
		if pr, _ := repo.GetPullRequest(context.Background(), id); pr != nil {
			prs = append(prs, pr)
		}
	}

	return prs, nil
}

// SavePullRequest сохраняет PR в БД
func (repo *PostgresPullRequestRepository) SavePullRequest(ctx context.Context, pr *entity.PullRequest) error {
	query := `
		INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, pr_status, created_at, merged_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := repo.db.Pool.Exec(ctx, query,
		pr.PullRequestID,
		pr.PullRequestName,
		pr.AuthorID, pr.Status,
		pr.CreatedAt,
		pr.MergedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save pull request: %w", err)
	}

	return nil
}

// UpdatePullRequest выполняет запрос к БД для обновления
// изменяемой информации о PR
func (repo *PostgresPullRequestRepository) UpdatePullRequest(ctx context.Context, pr *entity.PullRequest) error {
	query := `
		UPDATE pull_requests 
		SET pull_request_name = $1, author_id = $2, pr_status = $3, created_at = $4, merged_at = $5
		WHERE pull_request_id = $6;
	`

	result, err := repo.db.Pool.Exec(ctx, query,
		pr.PullRequestName,
		pr.AuthorID,
		string(pr.Status),
		pr.CreatedAt,
		pr.MergedAt,
		pr.PullRequestID,
	)

	if err != nil {
		return fmt.Errorf("failed to update pull request: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("pr not found")
	}

	return nil
}

// GetOpenedPullRequestCount выполняет запрос к БД для
// получения числа PR в статусе OPEN
func (repo *PostgresPullRequestRepository) GetOpenedPullRequestCount(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM pull_requests
		WHERE pr_status = 'OPEN'
	`

	var count int
	err := repo.db.Pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get open PR count: %w", err)
	}

	return count, nil
}

// GetMergedPullRequestCount выполняет запрос к БД для
// получения числа PR в статусе MERGED
func (repo *PostgresPullRequestRepository) GetMergedPullRequestCount(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM pull_requests
		WHERE pr_status = 'MERGED'
	`

	var count int
	err := repo.db.Pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get merged PR count: %w", err)
	}

	return count, nil
}
