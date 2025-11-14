package user

import "github.com/salex06/pr-service/internal/model"

type UserRepository interface {
	GetUser(userId string) *model.User
	UpdateUser(user *model.User) *model.User
	SaveUser(user *model.User) *model.User
	UserExists(userId string) bool

	ChooseReviewers(prAuthor *model.User) []string
	ReassignReviewer(teamName string, idsExclusionList []string) *string
}
