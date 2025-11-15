package pr

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/salex06/pr-service/internal/database"
	"github.com/salex06/pr-service/internal/model"
)

type PostgresPullRequestRepository struct {
	db *database.DB
}

func NewPostgresPullRequestRepository(db *database.DB) *PostgresPullRequestRepository {
	return &PostgresPullRequestRepository{db: db}
}

func (repo *PostgresPullRequestRepository) PullRequestExists(ctx context.Context, prId string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1)
	`

	err := repo.db.Pool.QueryRow(ctx, query, prId).Scan(&exists)

	return exists, err
}

func (repo *PostgresPullRequestRepository) GetPullRequest(ctx context.Context, prId string) (*model.PullRequest, error) {
	query := `
		SELECT pull_request_id, pull_request_name, author_id, pr_status, created_at, merged_at 
		FROM pull_requests
		WHERE pull_request_id = $1
	`

	var pr model.PullRequest
	err := repo.db.Pool.QueryRow(ctx, query, prId).Scan(
		&pr.PullRequestId,
		&pr.PullRequestName,
		&pr.AuthorId,
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

func (repo *PostgresPullRequestRepository) GetPullRequests(ctx context.Context, prIds []string) ([]*model.PullRequest, error) {
	prs := make([]*model.PullRequest, 0, len(prIds))
	for _, id := range prIds {
		if pr, _ := repo.GetPullRequest(context.Background(), id); pr != nil {
			prs = append(prs, pr)
		}
	}

	return prs, nil
}

func (repo *PostgresPullRequestRepository) SavePullRequest(ctx context.Context, pr *model.PullRequest) error {
	query := `
		INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, pr_status, created_at, merged_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := repo.db.Pool.Exec(ctx, query,
		pr.PullRequestId,
		pr.PullRequestName,
		pr.AuthorId, pr.Status,
		pr.CreatedAt,
		pr.MergedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save pull request: %w", err)
	}

	return nil
}

func (repo *PostgresPullRequestRepository) UpdatePullRequest(ctx context.Context, pr *model.PullRequest) error {
	query := `
		UPDATE pull_requests 
		SET pull_request_name = $1, author_id = $2, pr_status = $3, created_at = $4, merged_at = $5
		WHERE pull_request_id = $6;
	`

	result, err := repo.db.Pool.Exec(ctx, query,
		pr.PullRequestName,
		pr.AuthorId,
		string(pr.Status),
		pr.CreatedAt,
		pr.MergedAt,
		pr.PullRequestId,
	)

	if err != nil {
		return fmt.Errorf("failed to update pull request: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("pr not found")
	}

	return nil
}
