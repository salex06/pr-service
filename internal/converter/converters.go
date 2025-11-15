package converter

import (
	"github.com/salex06/pr-service/internal/dto"
	"github.com/salex06/pr-service/internal/model"
)

func ConvertTeamMemberToUser(member *dto.TeamMember, teamName string) *model.User {
	return &model.User{
		UserId:   member.UserId,
		Username: member.Username,
		TeamName: teamName,
		IsActive: member.IsActive,
	}
}

func ConvertUserToTeamMember(user *model.User) *dto.TeamMember {
	return &dto.TeamMember{
		UserId:   user.UserId,
		Username: user.Username,
		IsActive: user.IsActive,
	}
}

func ConvertUsersToTeamMembers(users []*model.User) []*dto.TeamMember {
	converted := make([]*dto.TeamMember, 0, len(users))
	for _, user := range users {
		converted = append(converted, ConvertUserToTeamMember(user))
	}

	return converted
}

func ConvertUserModelToDto(user *model.User) *dto.User {
	return &dto.User{
		UserId:   user.UserId,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
}

func ConvertDtoToPr(pr *dto.PullRequest) *model.PullRequest {
	return &model.PullRequest{
		PullRequestId:   pr.PullRequestId,
		PullRequestName: pr.PullRequestName,
		AuthorId:        pr.AuthorId,
		Status:          pr.Status,
		CreatedAt:       pr.CreatedAt,
		MergedAt:        pr.MergedAt,
	}
}

func ConvertPrToDto(pr *model.PullRequest, reviewers []string) *dto.PullRequest {
	return &dto.PullRequest{
		PullRequestId:     pr.PullRequestId,
		PullRequestName:   pr.PullRequestName,
		AuthorId:          pr.AuthorId,
		Status:            pr.Status,
		AssignedReviewers: reviewers,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}

func ConvertPrToReassigningDto(pr *model.PullRequest, reviewers []string, replacedBy string) *dto.ReassignPrResponse {
	return &dto.ReassignPrResponse{
		Pr: dto.PullRequest{
			PullRequestId:     pr.PullRequestId,
			PullRequestName:   pr.PullRequestName,
			AuthorId:          pr.AuthorId,
			Status:            pr.Status,
			AssignedReviewers: reviewers,
			CreatedAt:         pr.CreatedAt,
			MergedAt:          pr.MergedAt,
		},
		ReplacedBy: replacedBy,
	}
}

func ConvertPRsToAssignedPRs(userId string, prs []*model.PullRequest) *dto.AssignedPullRequests {
	pullRequestsShort := ConvertPrToShortPr(prs)

	return &dto.AssignedPullRequests{
		UserId:       userId,
		PullRequests: pullRequestsShort,
	}
}

func ConvertPrToShortPr(prs []*model.PullRequest) []dto.PullRequestShort {
	converted := make([]dto.PullRequestShort, 0, len(prs))

	for _, v := range prs {
		converted = append(converted, dto.PullRequestShort{
			PullRequestId:   v.PullRequestId,
			PullRequestName: v.PullRequestName,
			AuthorId:        v.AuthorId,
			Status:          v.Status,
		})
	}

	return converted
}
