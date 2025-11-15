package dto

// AppStat - основная структура для хранения и передачи статистики сервиса
type AppStat struct {
	TotalUsersCount  int `json:"total_users_count"`
	ActiveUsersCount int `json:"active_users_count"`

	TotalTeamsCount int `json:"total_teams_count"`

	OpenedPRCount int `json:"opened_pull_requests_count"`
	MergedPRCount int `json:"merged_pull_requests_count"`

	UserCountByTeam []*TeamSize `json:"users_count_by_team"`

	AssignmentsCountByUser []*AssignmentsByUser `json:"assignments_count_by_user"`
}

// TeamSize представляет структуру для хранения
// пар "название команды - количество членов команды"
type TeamSize struct {
	TeamName  string `json:"team_name"`
	UserCount int    `json:"user_count"`
}

// AssignmentsByUser представляет структуру для хранения
// пар "идентификатор пользователя - количество назначений данного пользователя на PR"
type AssignmentsByUser struct {
	UserID           string `json:"user_id"`
	AssignmentsCount int    `json:"assignments_count"`
}
