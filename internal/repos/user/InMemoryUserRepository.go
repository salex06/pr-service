package user

import (
	"context"
	"fmt"
	"slices"

	"github.com/salex06/pr-service/internal/model"
)

type InMemoryUserRepository struct {
	storage map[string]*model.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		storage: make(map[string]*model.User),
	}
}

func (db *InMemoryUserRepository) GetUser(ctx context.Context, userId string) (*model.User, error) {
	if user, ok := db.storage[userId]; ok {
		return user, nil
	}

	return nil, fmt.Errorf("user not found")
}

func (db *InMemoryUserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	db.storage[user.UserId] = user

	return nil
}

func (db *InMemoryUserRepository) SaveUser(ctx context.Context, user *model.User) error {
	db.storage[user.UserId] = user

	return nil
}

func (db *InMemoryUserRepository) UserExists(context context.Context, userId string) (bool, error) {
	_, ok := db.storage[userId]

	return ok, nil
}

func (db *InMemoryUserRepository) ChooseReviewers(context context.Context, prAuthor *model.User) ([]string, error) {
	activeTeammates := make([]*model.User, 0)
	for _, v := range db.storage {
		if v.IsActive && v.TeamName == prAuthor.TeamName && v.UserId != prAuthor.UserId {
			activeTeammates = append(activeTeammates, v)
		}
	}

	//Пока что выбираются двое первых сотрудников
	switch len(activeTeammates) {
	case 0:
		return make([]string, 0), nil
	case 1:
		return []string{activeTeammates[0].UserId}, nil
	default:
		return []string{
			activeTeammates[0].UserId,
			activeTeammates[1].UserId,
		}, nil
	}
}

func (db *InMemoryUserRepository) ReassignReviewer(ctx context.Context, teamName string, exclusionList []string) (*string, error) {
	for _, v := range db.storage {
		if v.IsActive && !slices.Contains(exclusionList, v.UserId) && v.TeamName == teamName {
			return &v.UserId, nil
		}
	}

	return nil, nil
}

func (db *InMemoryUserRepository) GetTeamMembers(ctx context.Context, teamName string) ([]*model.User, error) {
	members := make([]*model.User, 0)
	for _, v := range db.storage {
		if v.TeamName == teamName {
			members = append(members, v)
		}
	}
	return members, nil
}
