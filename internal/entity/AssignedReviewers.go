package entity

// AssignedReviewers представляет собой сущность,
// связывающую PR с назначенными сотрудниками
type AssignedReviewers struct {
	UserID        string
	PullRequestID string
}
