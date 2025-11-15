package user

import (
	"context"
	"fmt"
	"slices"

	"github.com/salex06/pr-service/internal/dto"
	"github.com/salex06/pr-service/internal/entity"
)

// InMemoryUserRepository представляет собой компонент,
// отвечающей за взаимодействие с in-memory БД (map),
// где хранится информация о пользователях
type InMemoryUserRepository struct {
	storage map[string]*entity.User
}

// NewInMemoryUserRepository конструирует и возвращает объект InMemoryUserRepository
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		storage: make(map[string]*entity.User),
	}
}

// GetUser возвращает пользователя с заданным userID (если не найден - nil)
func (db *InMemoryUserRepository) GetUser(ctx context.Context, userID string) (*entity.User, error) {
	if user, ok := db.storage[userID]; ok {
		return user, nil
	}

	return nil, fmt.Errorf("user not found")
}

// UpdateUser обновляет изменяемую информацию о пользователе
func (db *InMemoryUserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	db.storage[user.UserID] = user

	return nil
}

// SaveUser сохраняет пользователя в in-memory хранилище
func (db *InMemoryUserRepository) SaveUser(ctx context.Context, user *entity.User) error {
	db.storage[user.UserID] = user

	return nil
}

// UserExists проверяет, существует ли в in-memory хранилище
// пользователь с заданным id, и возвращает результат
func (db *InMemoryUserRepository) UserExists(ctx context.Context, userID string) (bool, error) {
	_, ok := db.storage[userID]

	return ok, nil
}

// ChooseReviewers возвращает до 2-х сотрудников (первых попавшихся),
// которые могут быть назначены на PR
func (db *InMemoryUserRepository) ChooseReviewers(ctx context.Context, prAuthor *entity.User) ([]string, error) {
	activeTeammates := make([]*entity.User, 0)
	for _, v := range db.storage {
		if v.IsActive && v.TeamName == prAuthor.TeamName && v.UserID != prAuthor.UserID {
			activeTeammates = append(activeTeammates, v)
		}
	}

	// Пока что выбираются двое первых сотрудников
	switch len(activeTeammates) {
	case 0:
		return make([]string, 0), nil
	case 1:
		return []string{activeTeammates[0].UserID}, nil
	default:
		return []string{
			activeTeammates[0].UserID,
			activeTeammates[1].UserID,
		}, nil
	}
}

// ReassignReviewer выполняет запрос к БД и возвращает
// первого попавшегося активного сотрудника из команды автора PR,
// который может быть назначен на PR вместо прежнего сотрудника
func (db *InMemoryUserRepository) ReassignReviewer(ctx context.Context, teamName string, exclusionList []string) (*string, error) {
	for _, v := range db.storage {
		if v.IsActive && !slices.Contains(exclusionList, v.UserID) && v.TeamName == teamName {
			return &v.UserID, nil
		}
	}

	return nil, nil
}

// GetTeamMembers возвращает сотрудников,
// которые являются членами заданной команды
func (db *InMemoryUserRepository) GetTeamMembers(ctx context.Context, teamName string) ([]*entity.User, error) {
	members := make([]*entity.User, 0)
	for _, v := range db.storage {
		if v.TeamName == teamName {
			members = append(members, v)
		}
	}
	return members, nil
}

// GetTotalUserCount возвращает общее количество пользователей
func (db *InMemoryUserRepository) GetTotalUserCount(ctx context.Context) (int, error) {
	return len(db.storage), nil
}

// GetActiveUserCount возвращает количество активных пользователей
func (db *InMemoryUserRepository) GetActiveUserCount(ctx context.Context) (int, error) {
	count := 0

	for _, v := range db.storage {
		if v.IsActive {
			count++
		}
	}

	return count, nil
}

// GetUserCountByTeam выполняет запрос к БД и возвращает
// количество пользователей в каждой команде
func (db *InMemoryUserRepository) GetUserCountByTeam(ctx context.Context) ([]*dto.TeamSize, error) {
	temp := make(map[string]int, 0)
	for _, v := range db.storage {
		temp[v.TeamName]++
	}

	teamSizes := make([]*dto.TeamSize, 0, len(temp))
	for k, v := range temp {
		teamSizes = append(teamSizes, &dto.TeamSize{
			TeamName:  k,
			UserCount: v,
		})
	}

	return teamSizes, nil
}
