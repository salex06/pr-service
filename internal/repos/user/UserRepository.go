// Package user - пакет с репозиториями, отвечающими за взаимодействие с БД,
// где хранится информация о пользователях
package user

import (
	"context"

	"github.com/salex06/pr-service/internal/entity"
)

// UserRepository представляет интерфейс взаимодействия
// с базой данных, где хранится информация о пользователях
type UserRepository interface {
	GetUser(ctx context.Context, userID string) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	SaveUser(ctx context.Context, user *entity.User) error
	UserExists(ctx context.Context, userID string) (bool, error)

	GetTeamMembers(ctx context.Context, teamName string) ([]*entity.User, error)

	ChooseReviewers(ctx context.Context, prAuthor *entity.User) ([]string, error)
	ReassignReviewer(ctx context.Context, teamName string, idsExclusionList []string) (*string, error)
}
