// Package converter - пакет, определяющий функции для преобразования одних структур в другие
package converter

import (
	"github.com/salex06/pr-service/internal/dto"
	"github.com/salex06/pr-service/internal/entity"
)

// ConvertTeamMemberToUser преобразовывает форму представления TeamMember
// и название команды в сущность User
func ConvertTeamMemberToUser(member *dto.TeamMember, teamName string) *entity.User {
	return &entity.User{
		UserID:   member.UserID,
		Username: member.Username,
		TeamName: teamName,
		IsActive: member.IsActive,
	}
}

// ConvertUserToTeamMember преобразовывает сущность User
// в форму представления сущности - TeamMember
func ConvertUserToTeamMember(user *entity.User) *dto.TeamMember {
	return &dto.TeamMember{
		UserID:   user.UserID,
		Username: user.Username,
		IsActive: user.IsActive,
	}
}

// ConvertUsersToTeamMembers преобразовывает слайс сущностей User
// в слайс объектов формы представления TeamMember
func ConvertUsersToTeamMembers(users []*entity.User) []*dto.TeamMember {
	converted := make([]*dto.TeamMember, 0, len(users))
	for _, user := range users {
		converted = append(converted, ConvertUserToTeamMember(user))
	}

	return converted
}

// ConvertUserEntityToDto преобразовывает сущность User
// в форму представления User
func ConvertUserEntityToDto(user *entity.User) *dto.User {
	return &dto.User{
		UserID:   user.UserID,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
}

// ConvertPrDtoToPrEntity преобразовывает форму представления PullRequest
// в сущность PullRequest
func ConvertPrDtoToPrEntity(pr *dto.PullRequest) *entity.PullRequest {
	return &entity.PullRequest{
		PullRequestID:   pr.PullRequestID,
		PullRequestName: pr.PullRequestName,
		AuthorID:        pr.AuthorID,
		Status:          pr.Status,
		CreatedAt:       pr.CreatedAt,
		MergedAt:        pr.MergedAt,
	}
}

// ConvertPrToDto преобразовывает сущность PullRequest и список
// назначенных ревьюеров в форму представления PullRequest
func ConvertPrToDto(pr *entity.PullRequest, reviewers []string) *dto.PullRequest {
	return &dto.PullRequest{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            pr.Status,
		AssignedReviewers: reviewers,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}

// ConvertPrToReassigningDto преобразовывает сущность PR,
// назначенных ревьюеров и идентификатором вновь назначенного сотрудника
// в структуру ReassignPrResponce
func ConvertPrToReassigningDto(pr *entity.PullRequest, reviewers []string, replacedBy string) *dto.ReassignPrResponse {
	return &dto.ReassignPrResponse{
		Pr: dto.PullRequest{
			PullRequestID:     pr.PullRequestID,
			PullRequestName:   pr.PullRequestName,
			AuthorID:          pr.AuthorID,
			Status:            pr.Status,
			AssignedReviewers: reviewers,
			CreatedAt:         pr.CreatedAt,
			MergedAt:          pr.MergedAt,
		},
		ReplacedBy: replacedBy,
	}
}

// ConvertPRsToAssignedPRs преобразовывает набор сущностей PullRequest,
// на которые назначен сотрудник с идентификатором userID, в структуру
// AssignedPullRequests
func ConvertPRsToAssignedPRs(userID string, prs []*entity.PullRequest) *dto.AssignedPullRequests {
	pullRequestsShort := ConvertPrToShortPr(prs)

	return &dto.AssignedPullRequests{
		UserID:       userID,
		PullRequests: pullRequestsShort,
	}
}

// ConvertPrToShortPr преобразовывает сущность PullRequest
// в его краткую форму представления PullRequestShort
func ConvertPrToShortPr(prs []*entity.PullRequest) []dto.PullRequestShort {
	converted := make([]dto.PullRequestShort, 0, len(prs))

	for _, v := range prs {
		converted = append(converted, dto.PullRequestShort{
			PullRequestID:   v.PullRequestID,
			PullRequestName: v.PullRequestName,
			AuthorID:        v.AuthorID,
			Status:          v.Status,
		})
	}

	return converted
}
