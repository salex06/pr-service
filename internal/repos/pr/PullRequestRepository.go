package pr

import "github.com/salex06/pr-service/internal/model"

type PullRequestRepository interface {
	PullRequestExists(prId string) bool

	GetPullRequest(prId string) *model.PullRequest
	GetPullRequests(prIds []string) []*model.PullRequest

	SavePullRequest(pr *model.PullRequest) *model.PullRequest
	UpdatePullRequest(pr *model.PullRequest) *model.PullRequest
}
