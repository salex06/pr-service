package model

import "time"

type PullRequestStatus string

const (
	OPEN   PullRequestStatus = "OPEN"
	MERGED PullRequestStatus = "MERGED"
)

type PullRequest struct {
	PullRequestId   string
	PullRequestName string
	AuthorId        string
	Status          PullRequestStatus
	CreatedAt       *time.Time
	MergedAt        *time.Time
}
