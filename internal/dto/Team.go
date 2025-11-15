package dto

// Team является формой представления сущности Team
// с названием команды и её представителями
type Team struct {
	TeamName string        `json:"team_name"`
	Members  []*TeamMember `json:"members"`
}
