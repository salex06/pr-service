package user

import (
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

func (db *InMemoryUserRepository) GetUser(userId string) *model.User {
	if user, ok := db.storage[userId]; ok {
		return user
	}
	return nil
}

func (db *InMemoryUserRepository) UpdateUser(user *model.User) *model.User {
	db.storage[user.UserId] = user

	return db.storage[user.UserId]
}

func (db *InMemoryUserRepository) SaveUser(user *model.User) *model.User {
	db.storage[user.UserId] = user

	return db.storage[user.UserId]
}

func (db *InMemoryUserRepository) UserExists(userId string) bool {
	_, ok := db.storage[userId]

	return ok
}

func (db *InMemoryUserRepository) ChooseReviewers(prAuthor *model.User) []string {
	activeTeammates := make([]*model.User, 0)
	for _, v := range db.storage {
		if v.IsActive && v.TeamName == prAuthor.TeamName && v.UserId != prAuthor.UserId {
			activeTeammates = append(activeTeammates, v)
		}
	}

	//Пока что выбираются двое первых сотрудников
	switch len(activeTeammates) {
	case 0:
		return make([]string, 0)
	case 1:
		return []string{activeTeammates[0].UserId}
	default:
		return []string{
			activeTeammates[0].UserId,
			activeTeammates[1].UserId,
		}
	}
}

func (db *InMemoryUserRepository) ReassignReviewer(teamName string, exclusionList []string) *string {
	for _, v := range db.storage {
		if v.IsActive && !slices.Contains(exclusionList, v.UserId) && v.TeamName == teamName {
			return &v.UserId
		}
	}

	return nil
}
