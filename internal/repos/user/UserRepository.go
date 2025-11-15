package user

import (
	"context"

	"github.com/salex06/pr-service/internal/model"
)

type UserRepository interface {
	GetUser(ctx context.Context, userId string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	SaveUser(ctx context.Context, user *model.User) error
	UserExists(ctx context.Context, userId string) (bool, error)

	GetTeamMembers(ctx context.Context, teamName string) ([]*model.User, error)

	ChooseReviewers(ctx context.Context, prAuthor *model.User) ([]string, error)
	ReassignReviewer(ctx context.Context, teamName string, idsExclusionList []string) (*string, error)
}
