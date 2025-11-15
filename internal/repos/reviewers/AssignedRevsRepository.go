package reviewers

import "context"

type AssignedRevsRepository interface {
	GetAssignedPullRequestIds(ctx context.Context, userId string) ([]string, error)
	GetAssignedReviewersIds(ctx context.Context, pullRequestId string) ([]string, error)

	CreateAssignment(ctx context.Context, userId string, prId string) error
	DeleteAssignment(ctx context.Context, userId string, prId string) error
}
