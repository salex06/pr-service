package dto

// TeamMember является формой представления сущности User
// с уникальным идентификатором, именем и флагом активности
type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}
