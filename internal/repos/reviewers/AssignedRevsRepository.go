package reviewers

type AssignedRevsRepository interface {
	GetAssignedPullRequestIds(userId string) []string
	GetAssignedReviewersIds(pullRequestId string) []string

	CreateAssignment(userId string, prId string)
	DeleteAssignment(userId string, prId string)
}
